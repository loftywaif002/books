package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/essentialbooks/books/pkg/common"
	"github.com/kjk/notionapi"
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
	FileNameBase string // TODO: possibly not needed

	Title     string // "Go", "jQuery" etcc
	titleSafe string
	TitleLong string // "Essential Go", "Essential jQuery" etc.

	NotionStartPageID string
	pageIDToPage      map[string]*notionapi.Page
	RootPage          *Page

	idToPage map[string]*Page

	fileCache      *FileCache
	SourceDir      string // dir where source markdown files are
	destDir        string // dir where destitation html files are
	SoContributors []SoContributor

	cachedArticlesCount int

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

// ContributorCount returns number of contributors
func (b *Book) ContributorCount() int {
	return len(b.SoContributors)
}

// ContributorsURL returns url of the chapter that lists contributors
func (b *Book) ContributorsURL() string {
	return b.URL() + "/contributors"
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
	panic("NYI")
	n := 0
	/*
		for _, ch := range b.Chapters {
			n += len(ch.Articles)
		}
		// each chapter has 000-index.md which is also an article
		n += len(b.Chapters)
	*/
	b.cachedArticlesCount = n
	return n
}

// ChaptersCount returns number of chapters
func (b *Book) ChaptersCount() int {
	panic("NYI")
	return 0
	// return len(b.Chapters)
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
