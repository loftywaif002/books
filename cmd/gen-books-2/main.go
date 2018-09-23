package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"html/template"

	"github.com/essentialbooks/books/pkg/common"
	"github.com/kjk/notionapi"
	"github.com/tdewolff/minify"
)

var (
	flgAnalytics string
	flgNoCache   bool

	soUserIDToNameMap map[int]string
	googleAnalytics   template.HTML
	doMinify          bool
	minifier          *minify.M
)

var (
	books = []*Book{
		&Book{
			Title:     "Go",
			TitleLong: "Essential Go",
			SourceDir: filepath.Join("books", "go"),
			// "https://www.notion.so/kjkpublic/Essential-Go-2cab1ed2b7a44584b56b0d3ca9b80185"
			NotionStartPageID: "2cab1ed2b7a44584b56b0d3ca9b80185",
		},
	}
)

const (
	// https://www.netlify.com/docs/headers-and-basic-auth/#custom-headers
	netlifyHeaders = `
# long-lived caching
/s/*
  Cache-Control: max-age=31536000
/*
  X-Content-Type-Options: nosniff
  X-Frame-Options: DENY
  X-XSS-Protection: 1; mode=block
`
)

const (
	googleAnalyticsTmpl = `<script async src="https://www.googletagmanager.com/gtag/js?id=%s"></script>
    <script>
        window.dataLayer = window.dataLayer || [];
        function gtag(){dataLayer.push(arguments);}
        gtag('js', new Date());
        gtag('config', '%s')
    </script>
`
)

func parseFlags() {
	flag.StringVar(&flgAnalytics, "analytics", "", "google analytics code")

	flag.BoolVar(&flgNoCache, "no-cache", false, "if true, disables cache for notion")
	flag.Parse()
}

func downloadBook(book *Book) {
	notionStartPageID := book.NotionStartPageID
	book.pageIDToPage = map[string]*notionapi.Page{}
	loadNotionPages(notionStartPageID, book.pageIDToPage, !flgNoCache)
	fmt.Printf("Loaded %d pages for book %s\n", len(book.pageIDToPage), book.Title)
	bookFromPages(book)
}

func iterPages(book *Book, onPage func(*Page) bool) {
	processed := map[string]bool{}
	toProcess := []*Page{book.RootPage}
	for len(toProcess) > 0 {
		page := toProcess[0]
		toProcess = toProcess[1:]
		id := normalizeID(page.NotionPage.ID)
		if processed[id] {
			continue
		}
		processed[id] = true
		toProcess = append(toProcess, page.Pages...)
		shouldContinue := onPage(page)
		if !shouldContinue {
			return
		}
	}
}

func buildIDToPage(book *Book) {
	book.idToPage = map[string]*Page{}
	fn := func(page *Page) bool {
		id := normalizeID(page.NotionPage.ID)
		book.idToPage[id] = page
		return true
	}
	iterPages(book, fn)
}

func bookPagesToHTML(book *Book) {
	nProcessed := 0
	fn := func(page *Page) bool {
		notionToHTML(page, book)
		nProcessed++
		return true
	}
	iterPages(book, fn)
	fmt.Printf("bookPagesToHTML: processed %d pages for book %s\n", nProcessed, book.TitleLong)
}

func genBookFiles(book *Book) {
	buildIDToPage(book)
	bookPagesToHTML(book)
}

func main() {
	parseFlags()

	//flgNoCache = true

	//maybeRemoveNotionCache()
	for _, book := range books {
		book.titleSafe = common.MakeURLSafe(book.Title)
		book.destDir = filepath.Join(destEssentialDir)

		downloadBook(book)
		genBookFiles(book)
	}
}
