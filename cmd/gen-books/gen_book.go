package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func gen404TopLevel() {
	d := struct {
		PageCommon
		Book *Book
	}{
		PageCommon: getPageCommon(),
	}
	path := filepath.Join(destDir, "404.html")
	execTemplateToFileMaybeMust("404.tmpl.html", d, path)
}

func genIndex(books []*Book) {
	d := struct {
		PageCommon
		Books []*Book
		//GitHubText string
		//GitHubURL  string
	}{
		PageCommon: getPageCommon(),
		Books:      books,
	}
	path := filepath.Join(destDir, "index.html")
	execTemplateToFileMaybeMust("index.tmpl.html", d, path)
}

func genIndexGrid(books []*Book) {
	d := struct {
		PageCommon
		Books []*Book
	}{
		PageCommon: getPageCommon(),
		Books:      books,
	}
	path := filepath.Join(destDir, "index-grid.html")
	execTemplateToFileMaybeMust("index-grid.tmpl.html", d, path)
}

func genFeedback() {
	d := getPageCommon()
	fmt.Printf("writing feedback.html\n")
	path := filepath.Join(destDir, "feedback.html")
	execTemplateToFileMaybeMust("feedback.tmpl.html", d, path)
}

func genAbout() {
	d := getPageCommon()
	fmt.Printf("writing about.html\n")
	path := filepath.Join(destDir, "about.html")
	execTemplateToFileMaybeMust("about.tmpl.html", d, path)
}

func genArticle(article *Article, currChapNo int) {
	addSitemapURL(article.CanonnicalURL())

	d := struct {
		PageCommon
		*Article
		CurrentChapterNo int
	}{
		PageCommon:       getPageCommon(),
		Article:          article,
		CurrentChapterNo: currChapNo,
	}

	path := article.destFilePath()
	execTemplateToFileSilentMaybeMust("article.tmpl.html", d, path)
}

func genChapter(chapter *Chapter, currNo int) {
	addSitemapURL(chapter.CanonnicalURL())
	for _, article := range chapter.Articles {
		genArticle(article, currNo)
	}

	path := chapter.destFilePath()
	d := struct {
		PageCommon
		*Chapter
		CurrentChapterNo int
	}{
		PageCommon:       getPageCommon(),
		Chapter:          chapter,
		CurrentChapterNo: currNo,
	}
	execTemplateToFileSilentMaybeMust("chapter.tmpl.html", d, path)

	for _, imagePath := range chapter.images {
		imageName := filepath.Base(imagePath)
		dst := chapter.destImagePath(imageName)
		copyFileMaybeMust(dst, imagePath)
	}
}

func genBook(book *Book) {
	fmt.Printf("Started genering book %s\n", book.Title)
	timeStart := time.Now()

	genBookTOCSearchMust(book)

	// generate index.html for the book
	err := os.MkdirAll(book.destDir, 0755)
	maybePanicIfErr(err)
	if err != nil {
		return
	}

	d := struct {
		PageCommon
		Book *Book
	}{
		PageCommon: getPageCommon(),
		Book:       book,
	}

	path := filepath.Join(book.destDir, "index.html")
	execTemplateToFileSilentMaybeMust("book_index.tmpl.html", d, path)

	path = filepath.Join(book.destDir, "404.html")
	execTemplateToFileSilentMaybeMust("404.tmpl.html", d, path)

	addSitemapURL(book.CanonnicalURL())

	for i, chapter := range book.Chapters {
		book.sem <- true
		book.wg.Add(1)
		go func(idx int, chap *Chapter) {
			genChapter(chap, idx)
			book.wg.Done()
			<-book.sem
		}(i+1, chapter)
	}
	book.wg.Wait()

	fmt.Printf("Generated %s, %d chapters, %d articles in %s\n", book.Title, len(book.Chapters), book.ArticlesCount(), time.Since(timeStart))
}
