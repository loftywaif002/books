package main

import (
	"fmt"
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

func makeRenderHookCodeBlock(defaultLang string) html.RenderNodeFunc {
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		codeBlock, ok := node.(*ast.CodeBlock)
		if !ok {
			return ast.GoToNext, false
		}
		lang := string(codeBlock.Info)
		if false {
			fmt.Printf("lang: '%s', code: %s\n", lang, string(codeBlock.Literal[:16]))
			io.WriteString(w, "\n<pre class=\"chroma\"><code>")
			//html.EscapeHTML(w, codeBlock.Literal)
			io.WriteString(w, "</code></pre>\n")
		} else {
			htmlHighlight(w, string(codeBlock.Literal), lang, defaultLang)
		}
		return ast.GoToNext, true
	}
}
