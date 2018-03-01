package main

import (
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/essentialbooks/books/pkg/kvstore"
)

var (
	// empty but not nil
	emptyStringSlice = make([]string, 0)
)

// Chapter represents a book chapter
type Chapter struct {
	// stable, globally unique (across all bookd) id
	// either imported Id from Stack Overflow or auto-generated by us
	// allows stable urls and being able to cross-reference articles
	ID         string
	Book       *Book
	ChapterDir string
	// full path to 000-index.md file
	indexFilePath string
	indexDoc      kvstore.Doc // content of 000-index.md file
	Title         string      // extracted from IndexKV, used in book_index.tmpl.html
	FileNameBase  string      // format: a-${ID}-${Title}, used for URL and .html file name
	Articles      []*Article
	No            int

	cachedHTML template.HTML

	// for search we extract headings from markdown source
	cachedHeadings []string

	// path for image files for this chapter in source directory
	images []string
}

// URL is used in book_index.tmpl.html
func (c *Chapter) URL() string {
	// /essential/go/a-4023-parsing-command-line-arguments-and-flags
	return fmt.Sprintf("/essential/%s/%s", c.Book.FileNameBase, c.FileNameBase)
}

// CanonnicalURL returns full url including host
func (c *Chapter) CanonnicalURL() string {
	return urlJoin(siteBaseURL, c.URL())
}

// GitHubText returns text we display in GitHub box
func (c *Chapter) GitHubText() string {
	return "Edit on GitHub"
}

// GitHubURL returns url to GitHub repo
func (c *Chapter) GitHubURL() string {
	return c.Book.GitHubURL() + "/" + c.ChapterDir
}

// GitHubEditURL returns url to edit 000-index.md document
func (c *Chapter) GitHubEditURL() string {
	bookDir := filepath.Base(c.Book.destDir)
	uri := gitHubBaseURL + "/blob/master/books/" + bookDir
	return uri + "/" + c.ChapterDir + "/000-index.md"
}

// GitHubIssueURL returns link for reporting an issue about an article on githbu
// https://github.com/essentialbooks/books/issues/new?title=${title}&body=${body}&labels=docs"
func (c *Chapter) GitHubIssueURL() string {
	title := fmt.Sprintf("Issue for chapter '%s'", c.Title)
	body := fmt.Sprintf("From URL: %s\nFile: %s\n", c.CanonnicalURL(), c.GitHubEditURL())
	return gitHubBaseURL + fmt.Sprintf("/issues/new?title=%s&body=%s&labels=docs", title, body)
}

func (c *Chapter) destFilePath() string {
	return filepath.Join(destEssentialDir, c.Book.FileNameBase, c.FileNameBase+".html")
}

func (c *Chapter) destImagePath(name string) string {
	return filepath.Join(destEssentialDir, c.Book.FileNameBase, name)
}

// HTML retruns html version of Body: field
func (c *Chapter) HTML() template.HTML {
	if c.cachedHTML != "" {
		return c.cachedHTML
	}
	s, err := c.indexDoc.GetValue("Body")
	if err != nil {
		return template.HTML("")
	}
	html := markdownToHTML([]byte(s), "", c.Book)
	c.cachedHTML = template.HTML(html)
	return c.cachedHTML
}

// Headings returns headings in markdown file
func (c *Chapter) Headings() []string {
	if c.cachedHeadings != nil {
		return c.cachedHeadings
	}
	s, err := c.indexDoc.GetValue("Body")
	if err != nil {
		return emptyStringSlice
	}
	headings := parseHeadingsFromMarkdown([]byte(s))
	if headings == nil {
		headings = emptyStringSlice
	}
	c.cachedHeadings = headings
	return headings
}

// TODO: get rid of IntroductionHTML, SyntaxHTML etc., convert to just Body in markdown format

// VersionsHTML returns html version of versions
func (c *Chapter) VersionsHTML() template.HTML {
	s, err := c.indexDoc.GetValue("VersionsHtml")
	if err != nil {
		s = ""
	}
	return template.HTML(s)
}

// IntroductionHTML retruns html version of Introduction:
func (c *Chapter) IntroductionHTML() template.HTML {
	s, err := c.indexDoc.GetValue("Introduction")
	if err != nil {
		return template.HTML("")
	}
	html := markdownToHTML([]byte(s), "", c.Book)
	return template.HTML(html)
}

// SyntaxHTML retruns html version of Syntax:
func (c *Chapter) SyntaxHTML() template.HTML {
	s, err := c.indexDoc.GetValue("Syntax")
	if err != nil {
		return template.HTML("")
	}
	html := markdownToHTML([]byte(s), "", c.Book)
	return template.HTML(html)
}

// RemarksHTML retruns html version of Remarks:
func (c *Chapter) RemarksHTML() template.HTML {
	s, err := c.indexDoc.GetValue("Remarks")
	if err != nil {
		return template.HTML("")
	}
	html := markdownToHTML([]byte(s), "", c.Book)
	return template.HTML(html)
}

// ContributorsHTML retruns html version of Contributors:
func (c *Chapter) ContributorsHTML() template.HTML {
	s, err := c.indexDoc.GetValue("Contributors")
	if err != nil {
		return template.HTML("")
	}
	html := markdownToHTML([]byte(s), "", c.Book)
	return template.HTML(html)
}
