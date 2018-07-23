package main

import "fmt"

func panicIfErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func panicIf(cond bool, format string, args ...interface{}) {
	if cond {
		s := fmt.Sprintf(format, args...)
		panic(s)
	}
}

var (
	books = []string{
		"Go", "Essential Go", notionGoStartPage,
	}
)

func downloadBook(bookShortName, bookName, notionStartPageID string) *Book {
	pages := loadNotionPages(notionStartPageID)
	fmt.Printf("Loaded %d pages for book %s\n", len(pages), bookName)
	book := bookFromPages(notionStartPageID, notionIDToPage)
	book.Title = bookShortName
	book.TitleLong = bookName
	return book
}

func genBookFiles(book *Book) {
	fmt.Printf("Generating files for book %s\n", book.TitleLong)
}

func main() {
	nBooks := len(books) / 3
	panicIf(nBooks*3 != len(books), "bad definition of books")
	maybeRemoveNotionCache()
	for i := 0; i < nBooks; i++ {
		book := downloadBook(books[i*3], books[i*3+1], books[i*3+2])
		genBookFiles(book)
	}
}
