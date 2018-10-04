package main

import (
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/kjk/notionapi"
)

// HeadingInfo describes header/sub header
type HeadingInfo struct {
	Text string
	ID   string
}

// Page represents a single page in a book
type Page struct {
	NotionPage *notionapi.Page
	Title      string
	// reference to parent page, nil if top-level page
	Parent *Page

	Book *Book

	No int
	// meta information extracted from page blocks
	NotionID string
	// for legacy pages this is an id. Might be used for redirects
	ID              string
	StackOverflowID string
	Search          []string // was SearchSynonyms

	// extracted from embed blocks
	SourceFiles []*EmbeddedSourceFile

	BodyHTML template.HTML

	// each page can contain sub-pages
	Pages []*Page

	ChapterDir string // TODO: no such thing anymore

	// filled during html generation
	Headings []HeadingInfo

	// TODO: those should come from notion_cache and downloaded during download
	// step to notion_cache
	images []string
}

// Siblings returns siblings of the page, to easily generate toc
func (p *Page) Siblings() []*Page {
	if p.Parent == nil {
		return nil
	}
	return p.Parent.Pages
}

// Body is a temporary alias for BodyHTML
func (p *Page) Body() template.HTML {
	return p.BodyHTML
}

// HTML is a temporary alias for BodyHTML
func (p *Page) HTML() template.HTML {
	return p.BodyHTML
}

// URL returns url of the page
func (p *Page) URL() string {
	book := p.Book
	bookTitle := book.Dir // should this be book.titleSafe ?
	id := p.NotionID
	title := urlify(p.Title)
	// /essentail/go/${id}-title
	return fmt.Sprintf("/essential/%s/%s-%s", bookTitle, id, title)
}

// CanonnicalURL returns full url including host
func (p *Page) CanonnicalURL() string {
	return urlJoin(siteBaseURL, p.URL())
}

// SuggestEditText returns text we display in GitHub box
func (p *Page) SuggestEditText() string {
	return "Suggest an edit"
}

// GitHubURL returns url to GitHub repo
func (p *Page) GitHubURL() string {
	return p.Book.GitHubURL() + "/" + p.ChapterDir
}

// SuggestEditURL returns url to edit 000-index.md document
func (p *Page) SuggestEditURL() string {
	return notionBaseURL + normalizeID(p.NotionID)
}

// GitHubIssueURL returns link for reporting an issue about an article on githbu
// https://github.com/essentialbooks/books/issues/new?title=${title}&body=${body}&labels=docs"
func (p *Page) GitHubIssueURL() string {
	title := fmt.Sprintf("Issue for chapter '%s'", p.Title)
	body := fmt.Sprintf("From URL: %s\nFile: %s\n", p.CanonnicalURL(), p.SuggestEditURL())
	return gitHubBaseURL + fmt.Sprintf("/issues/new?title=%s&body=%s&labels=docs", title, body)
}

func (p *Page) destFilePath() string {
	title := urlify(p.Title)
	fileName := p.NotionID + "-" + title + ".html"
	return filepath.Join(destEssentialDir, p.Book.Dir, fileName)
}

func (p *Page) destImagePath(name string) string {
	return filepath.Join(destEssentialDir, p.Book.Dir, name)
}

// PageTitle returns title for the page
// We want this to be unique for SEO purposes
func (p *Page) PageTitle() string {
	var a []string
	for p != nil {
		t := p.Title
		if t != "" {
			a = append(a, t)
		}
		p = p.Parent
	}
	reverseStringSlice(a)
	return strings.Join(a, " / ")
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

func bookPageFromNotionPage(book *Book, page *notionapi.Page, pageIDToPage map[string]*notionapi.Page) *Page {
	res := &Page{}
	res.NotionPage = page
	res.NotionID = normalizeID(page.ID)
	res.Title = page.Root.Title
	extractMeta(res)
	extractEmbeddedSourceFiles(res)
	subPages := getSubPages(page, pageIDToPage)

	// fmt.Printf("bookPageFromNotionPage: %s %s\n", normalizeID(page.ID), res.Meta.ID)

	for _, subPage := range subPages {
		bookPage := bookPageFromNotionPage(book, subPage, pageIDToPage)
		bookPage.Book = book
		res.Pages = append(res.Pages, bookPage)
	}
	return res
}

func bookFromPages(book *Book) {
	startPageID := book.NotionStartPageID
	page := book.pageIDToPage[startPageID]
	panicIf(page.Root.Type != notionapi.BlockPage, "start block is of type '%s' and not '%s'", page.Root.Type, notionapi.BlockPage)
	book.Title = page.Root.Title
	book.RootPage = bookPageFromNotionPage(book, page, book.pageIDToPage)
}
