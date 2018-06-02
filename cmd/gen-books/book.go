package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/kjk/u"
)

// SoContributor describes a StackOverflow contributor
type SoContributor struct {
	ID      int
	URLPart string
	Name    string
}

// Book represents a book
type Book struct {
	Title          string // "Go", "jQuery" etcc
	titleSafe      string
	TitleLong      string // "Essential Go", "Essential jQuery" etc.
	FileNameBase   string
	Chapters       []*Chapter
	sourceDir      string // dir where source markdown files are
	destDir        string // dif where destitation html files are
	SoContributors []SoContributor

	cachedArticlesCount int
	defaultLang         string // default programming language for programming examples
	knownUrls           []string

	// generated toc javascript data
	tocData []byte
	// url of combined tocData and app.js
	AppJSURL string

	// for concurrency
	sem chan bool
	wg  sync.WaitGroup
}

// ContributorCount returns number of contributors
func (b *Book) ContributorCount() int {
	return len(b.SoContributors)
}

// ContributorsURL returns url of the chapter that lists contributors
func (b *Book) ContributorsURL() string {
	return b.URL() + "/contributors"
}

// GitHubText returns text we show in GitHub link
func (b *Book) GitHubText() string {
	return "Edit on GitHub"
}

// GitHubURL returns link to GitHub for this book
func (b *Book) GitHubURL() string {
	return gitHubBaseURL + "/tree/master/books/" + filepath.Base(b.destDir)
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

// ArticlesCount returns total number of articles
func (b *Book) ArticlesCount() int {
	if b.cachedArticlesCount != 0 {
		return b.cachedArticlesCount
	}
	n := 0
	for _, ch := range b.Chapters {
		n += len(ch.Articles)
	}
	// each chapter has 000-index.md which is also an article
	n += len(b.Chapters)
	b.cachedArticlesCount = n
	return n
}

// ChaptersCount returns number of chapters
func (b *Book) ChaptersCount() int {
	return len(b.Chapters)
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

var didPrint = false

func printKnownURLS(a []string) {
	if didPrint {
		return
	}
	didPrint = true
	fmt.Printf("%d known urls\n", len(a))
	for _, s := range a {
		fmt.Printf("%s\n", s)
	}
}

// turn partial url like "20381" into a full url like "20381-installing"
func (b *Book) fixupURL(uri string) string {
	// skip uris that are not article/chapter uris
	if strings.Contains(uri, "/") {
		return uri
	}
	for _, known := range b.knownUrls {
		if uri == known {
			return uri
		}
		if strings.HasPrefix(known, uri) {
			//fmt.Printf("fixupURL: %s => %s\n", uri, known)
			return known
		}
	}
	fmt.Printf("fixupURL: didn't fix up: %s\n", uri)
	//printKnownURLS(knownURLS)
	return uri
}

func (b *Book) makeFixupURL() func(uri string) string {
	return func(uri string) string {
		return b.fixupURL(uri)
	}
}
