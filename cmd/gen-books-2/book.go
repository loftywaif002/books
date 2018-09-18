package main

import (
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/essentialbooks/books/pkg/common"
	"github.com/kjk/notionapi"
	"github.com/kjk/u"
)

// Book represents a book
type Book struct {
	Title     string // "Go", "jQuery" etcc
	titleSafe string
	TitleLong string // "Essential Go", "Essential jQuery" etc.
	RootPage  *Page

	destDir string // dir where destitation html files are
	//SoContributors []SoContributor

	defaultLang string // default programming language for programming examples
	knownUrls   []string

	// generated toc javascript data
	tocData []byte
	// url of combined tocData and app.js
	AppJSURL string

	idToPage map[string]*Page
}

// EmbeddedSourceFile represents source file present in the repository
// and embedded via https://www.onlinetool.io/gitoembed/
type EmbeddedSourceFile struct {
	EmbedURL string

	// name of the file
	FileName string
	// full path of the file
	Path string

	FileExists bool

	// content of the file after filtering
	Lines         []string
	cachedData    []byte
	cachedSha1Hex string
}

// Data returns content of the file
func (f *EmbeddedSourceFile) Data() []byte {
	if len(f.cachedData) == 0 {
		s := strings.Join(f.Lines, "\n")
		f.cachedData = []byte(s)
	}
	return f.cachedData
}

// RealSha1Hex returns hex version of sha1 of file content
func (f *EmbeddedSourceFile) RealSha1Hex() string {
	if f.cachedSha1Hex == "" {
		f.cachedSha1Hex = u.Sha1HexOfBytes(f.Data())
	}
	return f.cachedSha1Hex
}

func dataToLines(d []byte) []string {
	s := string(d)
	return strings.Split(s, "\n")
}

func readFilteredSourceFile(path string) ([]string, error) {
	d, err := common.ReadFileNormalized(path)
	if err != nil {
		return nil, err
	}
	lines := dataToLines(d)
	lines = removeAnnotationLines(lines)
	return lines, nil
}

// we don't want to show our // :show annotations in snippets
func removeAnnotationLines(lines []string) []string {
	var res []string
	prevWasEmpty := false
	for _, l := range lines {
		if strings.Contains(l, "// :show ") {
			continue
		}
		if len(l) == 0 && prevWasEmpty {
			continue
		}
		prevWasEmpty = len(l) == 0
		res = append(res, l)
	}
	return res
}

// Page represents a single page in a book
type Page struct {
	NotionPage *notionapi.Page
	Title      string
	// reference to parent page, nil if top-level page
	Parent *Page

	// meta information extracted from page blocks
	NotionID string
	// for legacy pages this is an id. Might be used for redirects
	ID              string
	StackOverflowID string
	Search          []string

	// extracted from embed blocks
	SourceFiles []*EmbeddedSourceFile

	BodyHTML template.HTML

	// each page can contain sub-pages
	Pages []*Page

	// to easily generate toc
	Siblings  []Page
	IsCurrent bool // only used when part of Siblings
}

// URL returns url of the page
func (p *Page) URL() string {
	return ""
}

func findSourceFileForEmbedURL(page *Page, uri string) *EmbeddedSourceFile {
	for _, f := range page.SourceFiles {
		if f.EmbedURL == uri {
			if f.FileExists {
				return f
			}
			return nil
		}
	}
	return nil
}

// extract sub page information and removes blocks that contain
// this info
func getSubPages(page *notionapi.Page, pageIDToPage map[string]*notionapi.Page) []*notionapi.Page {
	var res []*notionapi.Page
	toRemove := map[int]bool{}
	for idx, block := range page.Root.Content {
		if block.Type != notionapi.BlockPage {
			continue
		}
		toRemove[idx] = true
		id := normalizeID(block.ID)
		subPage := pageIDToPage[id]
		panicIf(subPage == nil, "no sub page for id %s", id)
		res = append(res, subPage)
	}
	removeBlocks(page, toRemove)
	return res
}

// MetaValue represents a single key: value meta-value
type MetaValue struct {
	Key   string
	Value string
}

// returns nil if this is not a meta-value block
// meta-value block is a plain text block in format:
// $key: value e.g. `$Id: 59`
func extractMetaValueFromBlock(block *notionapi.Block) *MetaValue {
	if block.Type != notionapi.BlockText {
		return nil
	}
	if len(block.InlineContent) != 1 {
		return nil
	}
	inline := block.InlineContent[0]
	// must be plain text
	if !inline.IsPlain() {
		return nil
	}

	// remove empty lines at the top
	s := strings.TrimSpace(inline.Text)
	if len(s) < 4 {
		return nil
	}
	if s[0] != '$' {
		return nil
	}
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return nil
	}
	key := strings.ToLower(strings.TrimSpace(parts[0]))
	value := strings.TrimSpace(parts[1])
	return &MetaValue{key, value}
}

// remove blocks whose indexes are in toRemove
func removeBlocks(page *notionapi.Page, toRemove map[int]bool) {
	if len(toRemove) == 0 {
		return
	}

	blocks := page.Root.Content
	n := 0
	for i, el := range blocks {
		if toRemove[i] {
			continue
		}
		blocks[n] = el
		n++
	}
	page.Root.Content = blocks[:n]

	ids := page.Root.ContentIDs
	n = 0
	for i, el := range ids {
		if toRemove[i] {
			continue
		}
		ids[n] = el
		n++
	}
	page.Root.ContentIDs = ids
}

// extracts PageMeta and updates Block.Content to remove the blocks that
// contained meta information
func extractMeta(p *Page) {
	page := p.NotionPage
	toRemove := map[int]bool{}
	for idx, block := range page.Root.Content {
		mv := extractMetaValueFromBlock(block)
		if mv == nil {
			continue
		}
		toRemove[idx] = true
		page.Root.Content[idx] = nil
		// fmt.Printf("'%s' = '%s'\n", mv.Key, mv.Value)
		switch mv.Key {
		case "$id":
			p.ID = mv.Value
		case "$soid":
			p.StackOverflowID = mv.Value
		case "$search":
			p.Search = strings.Split(mv.Value, ",")
			for i, s := range p.Search {
				p.Search[i] = strings.TrimSpace(s)
			}
		case "$score":
			// ignore
		default:
			panicIf(true, "unknown key '%s' in page with id %s", mv.Key, normalizeID(page.ID))
		}
	}
	removeBlocks(page, toRemove)
}

// https://www.onlinetool.io/gitoembed/widget?url=https%3A%2F%2Fgithub.com%2Fessentialbooks%2Fbooks%2Fblob%2Fmaster%2Fbooks%2Fgo%2F0020-basic-types%2Fbooleans.go
// to:
// books/go/0020-basic-types/booleans.go
// returns empty string if doesn't conform to what we expect
func gitoembedToRelativePath(uri string) string {
	parsed, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	switch parsed.Host {
	case "www.onlinetool.io", "onlinetool.io":
		// do nothing
	default:
		return ""
	}
	path := parsed.Path
	if path != "/gitoembed/widget" {
		return ""
	}
	uri = parsed.Query().Get("url")
	// https://github.com/essentialbooks/books/blob/master/books/go/0020-basic-types/booleans.go
	parsed, err = url.Parse(uri)
	if parsed.Host != "github.com" {
		return ""
	}
	path = strings.TrimPrefix(parsed.Path, "/essentialbooks/books/")
	if path == parsed.Path {
		return ""
	}
	// blob/master/books/go/0020-basic-types/booleans.go
	path = strings.TrimPrefix(path, "blob/")
	// master/books/go/0020-basic-types/booleans.go
	// those are branch names. Should I just strip first 2 elements from the path?
	path = strings.TrimPrefix(path, "master/")
	path = strings.TrimPrefix(path, "notion/")
	// books/go/0020-basic-types/booleans.go
	return path
}

func extractEmbeddedSourceFiles(p *Page) {
	wd, err := os.Getwd()
	panicIfErr(err)
	page := p.NotionPage
	for _, block := range page.Root.Content {
		if block.Type != notionapi.BlockEmbed {
			continue
		}
		uri := block.FormatEmbed.DisplaySource
		f := &EmbeddedSourceFile{
			EmbedURL: uri,
		}
		p.SourceFiles = append(p.SourceFiles, f)
		relativePath := gitoembedToRelativePath(uri)
		if relativePath == "" {
			fmt.Printf("Couldn't parse embed uri '%s'\n", uri)
			continue
		}
		// fmt.Printf("Embed uri: %s, relativePath: %s\n", uri, relativePath)
		f.FileName = filepath.Base(relativePath)
		f.Path = filepath.Join(wd, relativePath)
		f.Lines, err = readFilteredSourceFile(f.Path)
		if err != nil {
			fmt.Printf("Failed to read '%s' extracted from '%s', error: %s\n", f.Path, uri, err)
			continue
		}
		f.FileExists = true
	}
}

func bookPageFromNotionPage(page *notionapi.Page, pageIDToPage map[string]*notionapi.Page) *Page {
	res := &Page{}
	res.NotionPage = page
	res.Title = page.Root.Title
	extractMeta(res)
	extractEmbeddedSourceFiles(res)
	subPages := getSubPages(page, pageIDToPage)

	// fmt.Printf("bookPageFromNotionPage: %s %s\n", normalizeID(page.ID), res.Meta.ID)

	for _, subPage := range subPages {
		bookPage := bookPageFromNotionPage(subPage, pageIDToPage)
		res.Pages = append(res.Pages, bookPage)
	}
	return res
}

func bookFromPages(startPageID string, pageIDToPage map[string]*notionapi.Page) *Book {
	page := pageIDToPage[startPageID]
	panicIf(page.Root.Type != notionapi.BlockPage, "start block is of type '%s' and not '%s'", page.Root.Type, notionapi.BlockPage)
	book := &Book{}
	book.Title = page.Root.Title
	book.RootPage = bookPageFromNotionPage(page, pageIDToPage)

	return book
}
