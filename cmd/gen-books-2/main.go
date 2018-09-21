package main

import (
	"flag"
	"fmt"

	"html/template"

	"github.com/kjk/notionapi"
	"github.com/tdewolff/minify"
)

var (
	// "https://www.notion.so/kjkpublic/Essential-Go-2cab1ed2b7a44584b56b0d3ca9b80185"
	notionGoStartPage = "2cab1ed2b7a44584b56b0d3ca9b80185"

	flgAnalytics string
	flgNoCache   bool

	allBookDirs       []string
	soUserIDToNameMap map[int]string
	googleAnalytics   template.HTML
	doMinify          bool
	minifier          *minify.M
)

var (
	books = []string{
		"Go", "Essential Go", notionGoStartPage,
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

func downloadBook(bookShortName, bookName, notionStartPageID string) *Book {
	idToPage := map[string]*notionapi.Page{}
	loadNotionPages(notionGoStartPage, idToPage, !flgNoCache)
	fmt.Printf("Loaded %d pages for book %s\n", len(idToPage), bookName)
	book := bookFromPages(notionStartPageID, idToPage)
	book.Title = bookShortName
	book.TitleLong = bookName
	return book
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

	nBooks := len(books) / 3
	panicIf(nBooks*3 != len(books), "bad definition of books")
	//maybeRemoveNotionCache()
	for i := 0; i < nBooks; i++ {
		book := downloadBook(books[i*3], books[i*3+1], books[i*3+2])
		genBookFiles(book)
	}
}
