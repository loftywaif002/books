package main

import (
	"fmt"
	"io"
	"path"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

var (
	htmlFormatter  *html.Formatter
	highlightStyle *chroma.Style
)

// CodeBlockInfo represents info about code snippet
type CodeBlockInfo struct {
	Lang          string
	GitHubURI     string
	PlaygroundURI string
}

func init() {
	htmlFormatter = html.New(html.WithClasses(), html.TabWidth(2))
	panicIf(htmlFormatter == nil, "couldn't create html formatter")
	styleName := "monokailight"
	highlightStyle = styles.Get(styleName)
	panicIf(highlightStyle == nil, "didn't find style '%s'", styleName)

}

// gross hack: we need to change html generated by chroma
func fixupHTMLCodeBlock(htmlCode string, info *CodeBlockInfo) string {
	classLang := ""
	if info.Lang != "" {
		classLang = " lang-" + info.Lang
	}

	if info.GitHubURI == "" && info.PlaygroundURI == "" {
		html := fmt.Sprintf(`
<div class="code-box%s">
	<div>
		%s
	</div>
</div>`, classLang, htmlCode)
		return html
	}

	playgroundPart := ""
	if info.PlaygroundURI != "" {
		playgroundPart = fmt.Sprintf(`
<div class="code-box-playground">
	<a href="%s" target="_blank">try online</a>
</div>
`, info.PlaygroundURI)
	}

	gitHubPart := ""
	if info.GitHubURI != "" {
		// gitHubLoc is sth. like github.com/essentialbooks/books/books/go/main.go
		fileName := path.Base(info.GitHubURI)
		gitHubPart = fmt.Sprintf(`
<div class="code-box-github">
	<a href="%s" target="_blank">%s</a>
</div>`, info.GitHubURI, fileName)
	}

	html := fmt.Sprintf(`
<div class="code-box%s">
	<div>
	%s
	</div>
	<div class="code-box-nav">
		%s
		%s
	</div>
</div>`, classLang, htmlCode, playgroundPart, gitHubPart)
	return html
}

// based on https://github.com/alecthomas/chroma/blob/master/quick/quick.go
func htmlHighlight(w io.Writer, source, lang, defaultLang string) error {
	if lang == "" {
		lang = defaultLang
	}
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return htmlFormatter.Format(w, highlightStyle, it)
}
