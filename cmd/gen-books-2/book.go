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
	Pages     []*Page

	destDir string // dir where destitation html files are
	//SoContributors []SoContributor

	defaultLang string // default programming language for programming examples
	knownUrls   []string

	// generated toc javascript data
	tocData []byte
	// url of combined tocData and app.js
	AppJSURL string
}

// Page represents a singl
type Page struct {
	NotionPage *notionapi.Page
	Title      string
	// reference to parent page, nil if top-level page
	Parent *Page

	BodyHTML template.HTML

	// each page can contain sub-pages
	Pages []*Page

	Siblings  []Page
	IsCurrent bool // only used when part of Siblings
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

// PageMeta describe meta-information for a page
type PageMeta struct {
	// for legacy pages this is an id. Might be used for redirects
	ID              string
	StackOverflowID string
	NotionID        string
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
	s := block.Title
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
		switch mv.Key {
		case "id":
			res.ID = mv.Value
		case "soid":
			res.StackOverflowID = mv.Value

		default:
			panicIf(true, "unknown key '%s' in page with id %s", mv.Key, normalizeID(page.ID))
		}
	}
	// TODO: hacky because ContentIDs is not out of synca
	page.Root.Content = newBlocks
	return res
}

func bookPageFromNotionPage(page *notionapi.Page, pageIDToPage map[string]*notionapi.Page) *Page {
	res := &Page{}
	res.Title = page.Root.Title
	subPages := getSubPages(page, pageIDToPage)

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
	pages := getSubPages(page, pageIDToPage)
	for _, page := range pages {
		bookPage := bookPageFromNotionPage(page, pageIDToPage)
		book.Pages = append(book.Pages, bookPage)
	}
	return book
}
