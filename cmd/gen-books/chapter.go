package main

import (
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/essentialbooks/books/pkg/kvstore"
)

// Chapter represents a book chapter
type Chapter struct {
	*MarkdownFile

	Book       *Book
	ChapterDir string
	indexDoc   kvstore.Doc // content of 000-index.md file
	Articles   []*Article

	cachedHTML template.HTML

	// for search we extract headings from markdown source
	cachedHeadings []HeadingInfo

	// path for image files for this chapter in source directory
	images []string
}

// URL is used in book_index.tmpl.html
func (c *Chapter) URL() string {
	// /essential/go/4023-parsing-command-line-arguments-and-flags
	return fmt.Sprintf("/essential/%s/%s", c.Book.FileNameBase, c.FileNameBase)
}

// CanonnicalURL returns full url including host
func (c *Chapter) CanonnicalURL() string {
	return urlJoin(siteBaseURL, c.URL())
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
	s, err := c.indexDoc.Get("Body")
	if err != nil {
		return template.HTML("")
	}
	html := markdownToHTML([]byte(s), "", c.Book.makeFixupURL())
	c.cachedHTML = template.HTML(html)
	return c.cachedHTML
}

// Headings returns headings in markdown file
func (c *Chapter) Headings() []HeadingInfo {
	if c.cachedHeadings != nil {
		return c.cachedHeadings
	}
	s, err := c.indexDoc.Get("Body")
	if err != nil {
		return nil
	}
	headings := parseHeadingsFromMarkdown([]byte(s))
	c.cachedHeadings = headings
	return headings
}

// TODO: get rid of IntroductionHTML, SyntaxHTML etc., convert to just Body in markdown format

// VersionsHTML returns html version of versions
func (c *Chapter) VersionsHTML() template.HTML {
	s, err := c.indexDoc.Get("VersionsHtml")
	if err != nil {
		s = ""
	}
	return template.HTML(s)
}

// IntroductionHTML retruns html version of Introduction:
func (c *Chapter) IntroductionHTML() template.HTML {
	s, err := c.indexDoc.Get("Introduction")
	if err != nil {
		return template.HTML("")
	}
	html := markdownToHTML([]byte(s), "", c.Book.makeFixupURL())
	return template.HTML(html)
}

// SyntaxHTML retruns html version of Syntax:
func (c *Chapter) SyntaxHTML() template.HTML {
	s, err := c.indexDoc.Get("Syntax")
	if err != nil {
		return template.HTML("")
	}
	html := markdownToHTML([]byte(s), "", c.Book.makeFixupURL())
	return template.HTML(html)
}

// RemarksHTML retruns html version of Remarks:
func (c *Chapter) RemarksHTML() template.HTML {
	s, err := c.indexDoc.Get("Remarks")
	if err != nil {
		return template.HTML("")
	}
	html := markdownToHTML([]byte(s), "", c.Book.makeFixupURL())
	return template.HTML(html)
}

// ContributorsHTML retruns html version of Contributors:
func (c *Chapter) ContributorsHTML() template.HTML {
	s, err := c.indexDoc.Get("Contributors")
	if err != nil {
		return template.HTML("")
	}
	html := markdownToHTML([]byte(s), "", c.Book.makeFixupURL())
	return template.HTML(html)
}
