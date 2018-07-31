package main

import (
	"html/template"
	"strings"

	"github.com/kjk/notionapi"
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
}

// Page represents a single page in a book
type Page struct {
	NotionPage *notionapi.Page
	Title      string
	// reference to parent page, nil if top-level page
	Parent *Page
	Meta   *PageMeta

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

// PageMeta describe meta-information for a page
type PageMeta struct {
	NotionID string
	// for legacy pages this is an id. Might be used for redirects
	ID              string
	StackOverflowID string
	Search          []string
}

// extract sub page information and removes blocks that contain
// this info
func getSubPages(page *notionapi.Page, pageIDToPage map[string]*notionapi.Page) []*notionapi.Page {
	var res []*notionapi.Page
	var newBlocks []*notionapi.Block
	for _, block := range page.Root.Content {
		if block.Type != notionapi.BlockPage {
			newBlocks = append(newBlocks, block)
			continue
		}
		id := normalizeID(block.ID)
		subPage := pageIDToPage[id]
		panicIf(subPage == nil, "no sub page for id %s", id)
		res = append(res, subPage)
	}
	// this is a bit hacky as not ContentIDs is out of sync
	page.Root.Content = newBlocks
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

// extracts PageMeta and updates Block.Content to remove the blocks that
// contained meta information
func extractMeta(page *notionapi.Page) *PageMeta {
	res := &PageMeta{}
	var newBlocks []*notionapi.Block
	for _, block := range page.Root.Content {
		mv := extractMetaValueFromBlock(block)
		if mv == nil {
			newBlocks = append(newBlocks, block)
			continue
		}
		//fmt.Printf("'%s' = '%s'\n", mv.Key, mv.Value)
		switch mv.Key {
		case "$id":
			res.ID = mv.Value
		case "$soid":
			res.StackOverflowID = mv.Value
		case "$search":
			res.Search = strings.Split(mv.Value, ",")
			for i, s := range res.Search {
				res.Search[i] = strings.TrimSpace(s)
			}
		case "$score":
			// ignore
		default:
			panicIf(true, "unknown key '%s' in page with id %s", mv.Key, normalizeID(page.ID))
		}
	}
	// TODO: hacky because ContentIDs is now out of sync
	page.Root.Content = newBlocks
	return res
}

func bookPageFromNotionPage(page *notionapi.Page, pageIDToPage map[string]*notionapi.Page) *Page {
	res := &Page{}
	res.Title = page.Root.Title
	res.Meta = extractMeta(page)
	subPages := getSubPages(page, pageIDToPage)

	//fmt.Printf("bookPageFromNotionPage: %s %s\n", normalizeID(page.ID), res.Meta.ID)

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
