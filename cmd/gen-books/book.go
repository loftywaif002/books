package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/kjk/notionapi"
	"github.com/kjk/u"
)

// Book represents a book
type Book struct {
	//FileNameBase string // TODO: possibly not needed

	Title     string // "Essential Go", "Essential jQuery" etcc
	titleSafe string
	// TODO: remove TitleLong since Title is the same
	TitleLong string // "Essential Go", "Essential jQuery" etc.

	NotionStartPageID string
	pageIDToPage      map[string]*notionapi.Page
	RootPage          *Page   // a tree of pages
	cachedPages       []*Page // pages flattened into an array

	idToPage map[string]*Page

	Dir            string // directory name for the book e.g. "go"
	SoContributors []SoContributor

	defaultLang string // default programming language for programming examples
	knownUrls   []string

	// generated toc javascript data
	tocData []byte
	// url of combined tocData and app.js
	AppJSURL string

	// for concurrency
	sem chan bool
	wg  sync.WaitGroup
}

// SourceDir is where source files for a given book are
func (b *Book) SourceDir() string {
	return filepath.Join("books", b.Dir)
}

// this is where html etc. files for a book end up
func (b *Book) destDir() string {
	return filepath.Join(destEssentialDir, b.Dir)
}

// ContributorCount returns number of contributors
func (b *Book) ContributorCount() int {
	return len(b.SoContributors)
}

// ContributorsURL returns url of the chapter that lists contributors
func (b *Book) ContributorsURL() string {
	return b.URL() + "/9999-contributors"
}

// SuggestEditText returns text we show in GitHub link
func (b *Book) SuggestEditText() string {
	return "Suggest an edit"
}

// GitHubURL returns link to GitHub for this book
func (b *Book) GitHubURL() string {
	return gitHubBaseURL + "/tree/master/books/" + filepath.Base(b.destDir())
}

// URL returns url of the book, used in index.tmpl.html
func (b *Book) URL() string {
	return fmt.Sprintf("/essential/%s/", b.titleSafe)
}

// CanonnicalURL returns full url including host
func (b *Book) CanonnicalURL() string {
	return urlJoin(siteBaseURL, b.URL())
}

// ShareOnTwitterText returns text for sharing on twitter
func (b *Book) ShareOnTwitterText() string {
	return fmt.Sprintf(`"Essential %s" - a free programming book`, b.Title)
}

// CoverURL returns url to cover image
func (b *Book) CoverURL() string {
	coverName := langToCover[b.titleSafe]
	return fmt.Sprintf("/covers/%s.png", coverName)
}

// CoverFullURL returns a URL for the cover including host
func (b *Book) CoverFullURL() string {
	return urlJoin(siteBaseURL, b.CoverURL())
}

// CoverTwitterFullURL returns a URL for the cover including host
func (b *Book) CoverTwitterFullURL() string {
	coverName := langToCover[b.titleSafe]
	coverURL := fmt.Sprintf("/covers/twitter/%s.png", coverName)
	return urlJoin(siteBaseURL, coverURL)
}

// Chapters returns pages that are top-level chapters
func (b *Book) Chapters() []*Page {
	return b.RootPage.Pages
}

// GetAllPages returns all pages, flattened
func (b *Book) GetAllPages() []*Page {
	// to prevent infinite recursion if pages show up twice (shouldn't happen)
	if len(b.cachedPages) > 0 {
		return b.cachedPages
	}
	if b.RootPage == nil {
		return nil
	}
	seen := map[*Page]bool{}
	pages := []*Page{b.RootPage}
	curr := 0
	for curr < len(pages) {
		page := pages[curr]
		curr++
		if seen[page] {
			continue
		}
		seen[page] = true
		for i, p := range page.Pages {
			p.Parent = page
			p.No = i
		}
		pages = append(pages, page.Pages...)
	}
	return pages
}

// PagesCount returns total number of articles
func (b *Book) PagesCount() int {
	return len(b.GetAllPages()) - 1 // don't count top page
}

// ChaptersCount returns number of chapters
func (b *Book) ChaptersCount() int {
	return len(b.RootPage.Pages)
}

func updateBookAppJS(book *Book) {
	srcName := fmt.Sprintf("app-%s.js", book.titleSafe)
	path := filepath.Join("tmpl", "app.js")
	d, err := ioutil.ReadFile(path)
	maybePanicIfErr(err)
	if err != nil {
		return
	}
	if doMinify {
		d2, err := minifier.Bytes("text/javascript", d)
		maybePanicIfErr(err)
		if err == nil {
			fmt.Printf("Minified %s from %d => %d (saved %d)\n", srcName, len(d), len(d2), len(d)-len(d2))
			d = d2
		}
	}

	d = append(book.tocData, d...)
	sha1Hex := u.Sha1HexOfBytes(d)
	name := nameToSha1Name(srcName, sha1Hex)
	dst := filepath.Join("www", "s", name)
	err = ioutil.WriteFile(dst, d, 0644)
	maybePanicIfErr(err)
	if err != nil {
		return
	}
	book.AppJSURL = "/s/" + name
	fmt.Printf("Created %s\n", dst)
}
