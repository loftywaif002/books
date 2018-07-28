package main

import (
	"fmt"
	"io"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
)

var (
	htmlFormatter  *html.Formatter
	highlightStyle *chroma.Style
)

func init() {
	htmlFormatter = html.New(html.WithClasses(), html.TabWidth(2))
	panicIf(htmlFormatter == nil, "couldn't create html formatter")
	styleName := "monokailight"
	highlightStyle = styles.Get(styleName)
	panicIf(highlightStyle == nil, "didn't find style '%s'", styleName)

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

func makeRenderHookCodeBlock(defaultLang string) mdhtml.RenderNodeFunc {
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		codeBlock, ok := node.(*ast.CodeBlock)
		if !ok {
			return ast.GoToNext, false
		}
		lang := string(codeBlock.Info)
		if false {
			fmt.Printf("lang: '%s', code: %s\n", lang, string(codeBlock.Literal[:16]))
			io.WriteString(w, "\n<pre class=\"chroma\"><code>")
			mdhtml.EscapeHTML(w, codeBlock.Literal)
			io.WriteString(w, "</code></pre>\n")
		} else {
			htmlHighlight(w, string(codeBlock.Literal), lang, defaultLang)
		}
		return ast.GoToNext, true
	}
}
