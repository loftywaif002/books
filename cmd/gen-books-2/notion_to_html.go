package main

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"strconv"
	"strings"

	"github.com/alecthomas/template"
	"github.com/kjk/notionapi"
)

// HTMLGenerator is for notion -> HTML generation
type HTMLGenerator struct {
	f           *bytes.Buffer
	page        *notionapi.Page
	level       int
	nToggle     int
	err         error
	idToArticle func(string) *Page
}

// NewHTMLGenerator returns new HTMLGenerator
func NewHTMLGenerator(page *notionapi.Page) *HTMLGenerator {
	return &HTMLGenerator{
		f:    &bytes.Buffer{},
		page: page,
	}
}

// Gen returns generated HTML
func (g *HTMLGenerator) Gen() []byte {
	page := g.page.Root
	f := page.FormatPage
	g.writeString(`<p></p>`)
	if f != nil && f.PageFont == "mono" {
		g.writeString(`<div style="font-family: monospace">`)
	}
	g.genContent(g.page.Root)
	if f != nil && f.PageFont == "mono" {
		g.writeString(`</div>`)
	}
	return g.f.Bytes()
}

func isValidNotionID(id string) bool {
	// TODO: more strict i.e. check all characters are hex
	return len(id) == len("ea07db1b9bff415ab180b0525f3898f6")
}

// change https://www.notion.so/Advanced-web-spidering-with-Puppeteer-ea07db1b9bff415ab180b0525f3898f6
// =>
// /article/${id}
func (g *HTMLGenerator) maybeReplaceNotionLink(uri string) string {
	if !strings.HasPrefix(uri, "https://www.notion.so/") {
		return uri
	}
	parts := strings.Split(uri, "-")
	n := len(parts)
	if n < 2 {
		return uri
	}
	id := normalizeID(parts[n-1])
	if !isValidNotionID(id) {
		return uri
	}
	article := g.idToArticle(id)
	return article.URL()
}

func (g *HTMLGenerator) getURLAndTitleForBlock(block *notionapi.Block) (string, string) {
	id := normalizeID(block.ID)
	article := g.idToArticle(id)
	if article == nil {
		title := block.Title
		fmt.Printf("No article for id %s %s\n", id, title)
		url := "/article/" + id + "/" + urlify(title)
		return url, title
	}

	return article.URL(), article.Title
}

func (g *HTMLGenerator) genInlineBlock(b *notionapi.InlineBlock) {
	var start, close string
	if b.AttrFlags&notionapi.AttrBold != 0 {
		start += "<b>"
		close += "</b>"
	}
	if b.AttrFlags&notionapi.AttrItalic != 0 {
		start += "<i>"
		close += "</i>"
	}
	if b.AttrFlags&notionapi.AttrStrikeThrought != 0 {
		start += "<strike>"
		close += "</strike>"
	}
	if b.AttrFlags&notionapi.AttrCode != 0 {
		start += "<code>"
		close += "</code>"
	}
	skipText := false
	if b.Link != "" {
		link := g.maybeReplaceNotionLink(b.Link)
		start += fmt.Sprintf(`<a href="%s">%s</a>`, link, b.Text)
		skipText = true
	}
	if b.UserID != "" {
		start += fmt.Sprintf(`<span class="user">@%s</span>`, b.UserID)
		skipText = true
	}
	if b.Date != nil {
		// TODO: serialize date properly
		start += fmt.Sprintf(`<span class="date">@TODO: date</span>`)
		skipText = true
	}
	if !skipText {
		start += b.Text
	}
	g.writeString(start + close)
}

func (g *HTMLGenerator) getInline(blocks []*notionapi.InlineBlock) []byte {
	b := g.newBuffer()
	g.genInlineBlocks(blocks)
	return g.restoreBuffer(b)
}

func (g *HTMLGenerator) genInlineBlocks(blocks []*notionapi.InlineBlock) {
	for _, block := range blocks {
		g.genInlineBlock(block)
	}
}

func (g *HTMLGenerator) genBlockSurrouded(block *notionapi.Block, start, close string) {
	g.writeString(start + "\n")
	g.genInlineBlocks(block.InlineContent)
	g.level++
	g.genContent(block)
	g.level--
	g.writeString(close + "\n")
}

/*
v is expected to be
[
	[
		"foo"
	]
]
and we want to return "foo"
If not present or unexpected shape, return ""
is still visible
*/
func propsValueToText(v interface{}) string {
	if v == nil {
		return ""
	}

	// [ [ "foo" ]]
	a, ok := v.([]interface{})
	if !ok {
		return fmt.Sprintf("type1: %T", v)
	}
	// [ "foo" ]
	if len(a) == 0 {
		return ""
	}
	v = a[0]
	a, ok = v.([]interface{})
	if !ok {
		return fmt.Sprintf("type2: %T", v)
	}
	// "foo"
	if len(a) == 0 {
		return ""
	}
	v = a[0]
	str, ok := v.(string)
	if !ok {
		return fmt.Sprintf("type3: %T", v)
	}
	return str
}

func (g *HTMLGenerator) genCollectionView(block *notionapi.Block) {
	viewInfo := block.CollectionViews[0]
	view := viewInfo.CollectionView
	columns := view.Format.TableProperties
	s := `<table class="notion-table"><thead><tr>`
	for _, col := range columns {
		colName := col.Property
		colInfo := viewInfo.Collection.CollectionSchema[colName]
		name := colInfo.Name
		s += `<th>` + html.EscapeString(name) + `</th>`
	}
	s += `</tr></thead>`
	s += `<tbody>`
	for _, row := range viewInfo.CollectionRows {
		s += `<tr>`
		props := row.Properties
		for _, col := range columns {
			colName := col.Property
			v := props[colName]
			colVal := propsValueToText(v)
			if colVal == "" {
				// use &nbsp; so that empty row still shows up
				// could also set a min-height to 1em or sth. like that
				s += `<td>&nbsp;</td>`
			} else {
				//colInfo := viewInfo.Collection.CollectionSchema[colName]
				// TODO: format colVal according to colInfo
				s += `<td>` + html.EscapeString(colVal) + `</td>`
			}
		}
		s += `</tr>`
	}
	s += `</tbody>`
	s += `</table>`
	g.writeString(s)
}

// Children of BlockColumnList are BlockColumn blocks
func (g *HTMLGenerator) genColumnList(block *notionapi.Block) {
	panicIf(block.Type != notionapi.BlockColumnList, "unexpected block type '%s'", block.Type)
	nColumns := len(block.Content)
	panicIf(nColumns == 0, "has no columns")
	// TODO: for now equal width columns
	s := `<div class="column-list">`
	g.writeString(s)

	for _, col := range block.Content {
		// TODO: get column ration from col.FormatColumn.ColumnRation, which is float 0...1
		panicIf(col.Type != notionapi.BlockColumn, "unexpected block type '%s'", col.Type)
		g.writeString(`<div>`)
		g.genBlocks(col.Content)
		g.writeString(`</div>`)
	}

	s = `</div>`
	g.writeString(s)
}

func (g *HTMLGenerator) newBuffer() *bytes.Buffer {
	curr := g.f
	g.f = &bytes.Buffer{}
	return curr
}

func (g *HTMLGenerator) restoreBuffer(b *bytes.Buffer) []byte {
	d := g.f.Bytes()
	g.f = b
	return d
}

func (g *HTMLGenerator) genToggle(block *notionapi.Block) {
	panicIf(block.Type != notionapi.BlockToggle, "unexpected block type '%s'", block.Type)
	g.nToggle++
	id := strconv.Itoa(g.nToggle)

	inline := g.getInline(block.InlineContent)

	b := g.newBuffer()
	g.genBlocks(block.Content)
	inner := g.restoreBuffer(b)

	s := fmt.Sprintf(`<div style="width: 100%%; margin-top: 2px; margin-bottom: 1px;">
    <div style="display: flex; align-items: flex-start; width: 100%%; padding-left: 2px; color: rgb(66, 66, 65);">

        <div style="margin-right: 4px; width: 24px; flex-grow: 0; flex-shrink: 0; display: flex; align-items: center; justify-content: center; min-height: calc((1.5em + 3px) + 3px); padding-right: 2px;">
            <div id="toggle-toggle-%s" onclick="javascript:onToggleClick(this)" class="toggler" style="align-items: center; user-select: none; display: flex; width: 1.25rem; height: 1.25rem; justify-content: center; flex-shrink: 0;">

                <svg id="toggle-closer-%s" width="100%%" height="100%%" viewBox="0 0 100 100" style="fill: currentcolor; display: none; width: 0.6875em; height: 0.6875em; transition: transform 300ms ease-in-out; transform: rotateZ(180deg);">
                    <polygon points="5.9,88.2 50,11.8 94.1,88.2 "></polygon>
                </svg>

                <svg id="toggle-opener-%s" width="100%%" height="100%%" viewBox="0 0 100 100" style="fill: currentcolor; display: block; width: 0.6875em; height: 0.6875em; transition: transform 300ms ease-in-out; transform: rotateZ(90deg);">
                    <polygon points="5.9,88.2 50,11.8 94.1,88.2 "></polygon>
                </svg>
            </div>
        </div>

        <div style="flex: 1 1 0px; min-width: 1px;">
            <div style="display: flex;">
                <div style="padding-top: 3px; padding-bottom: 3px">%s</div>
            </div>

            <div style="margin-left: -2px; display: none" id="toggle-content-%s">
                <div style="display: flex; flex-direction: column;">
                    <div style="width: 100%%; margin-top: 2px; margin-bottom: 0px;">
                        <div style="color: rgb(66, 66, 65);">
							<div style="">
								%s
                                <!-- <div style="padding: 3px 2px;">text inside list</div> -->
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
`, id, id, id, string(inline), id, string(inner))
	g.writeString(s)
}

func (g *HTMLGenerator) writeString(s string) {
	io.WriteString(g.f, s)
}

func (g *HTMLGenerator) genBlock(block *notionapi.Block) {
	levelCls := ""
	if g.level > 0 {
		levelCls = fmt.Sprintf(" lvl%d", g.level)
	}

	switch block.Type {
	case notionapi.BlockText:
		start := fmt.Sprintf(`<p>`)
		close := `</p>`
		g.genBlockSurrouded(block, start, close)
	case notionapi.BlockHeader:
		start := fmt.Sprintf(`<h1 class="hdr%s">`, levelCls)
		close := `</h1>`
		g.genBlockSurrouded(block, start, close)
	case notionapi.BlockSubHeader:
		start := fmt.Sprintf(`<h2 class="hdr%s">`, levelCls)
		close := `</h2>`
		g.genBlockSurrouded(block, start, close)
	case notionapi.BlockTodo:
		clsChecked := ""
		if block.IsChecked {
			clsChecked = " todo-checked"
		}
		start := fmt.Sprintf(`<div class="todo%s%s">`, levelCls, clsChecked)
		close := `</div>`
		g.genBlockSurrouded(block, start, close)
	case notionapi.BlockToggle:
		g.genToggle(block)
	case notionapi.BlockQuote:
		start := fmt.Sprintf(`<blockquote class="%s">`, levelCls)
		close := `</blockquote>`
		g.genBlockSurrouded(block, start, close)
	case notionapi.BlockDivider:
		fmt.Fprintf(g.f, `<hr class="%s"/>`+"\n", levelCls)
	case notionapi.BlockPage:
		cls := "page"
		if block.IsLinkToPage() {
			cls = "page-link"
		}
		url, title := g.getURLAndTitleForBlock(block)
		title = template.HTMLEscapeString(title)
		html := fmt.Sprintf(`<div class="%s%s"><a href="%s">%s</a></div>`, cls, levelCls, url, title)
		fmt.Fprintf(g.f, "%s\n", html)
	case notionapi.BlockCode:
		/*
			code := template.HTMLEscapeString(block.Code)
			fmt.Fprintf(g.f, `<div class="%s">Lang for code: %s</div>
			<pre class="%s">
			%s
			</pre>`, levelCls, block.CodeLanguage, levelCls, code)
		*/
		htmlHighlight(g.f, string(block.Code), block.CodeLanguage, "")
	case notionapi.BlockBookmark:
		fmt.Fprintf(g.f, `<div class="bookmark %s">Bookmark to %s</div>`+"\n", levelCls, block.Link)
	case notionapi.BlockGist:
		s := fmt.Sprintf(`<script src="%s.js"></script>`, block.Source)
		g.writeString(s)
	case notionapi.BlockImage:
		link := block.ImageURL
		fmt.Fprintf(g.f, `<img class="%s" style="width: 100%%" src="%s" />`+"\n", levelCls, link)
	case notionapi.BlockColumnList:
		g.genColumnList(block)
	case notionapi.BlockCollectionView:
		g.genCollectionView(block)
	default:
		fmt.Printf("Unsupported block type '%s', id: %s\n", block.Type, block.ID)
		panic(fmt.Sprintf("Unsupported block type '%s'", block.Type))
	}
}

func (g *HTMLGenerator) genBlocks(blocks []*notionapi.Block) {
	for len(blocks) > 0 {
		block := blocks[0]
		if block == nil {
			fmt.Printf("Missing block\n")
			blocks = blocks[1:]
			continue
		}

		if block.Type == notionapi.BlockNumberedList {
			fmt.Fprintf(g.f, `<ol>`)
			for len(blocks) > 0 {
				block := blocks[0]
				if block.Type != notionapi.BlockNumberedList {
					break
				}
				g.genBlockSurrouded(block, `<li>`, `</li>`)
				blocks = blocks[1:]
			}
			fmt.Fprintf(g.f, `</ol>`)
		} else if block.Type == notionapi.BlockBulletedList {
			fmt.Fprintf(g.f, `<ul>`)
			for len(blocks) > 0 {
				block := blocks[0]
				if block.Type != notionapi.BlockBulletedList {
					break
				}
				g.genBlockSurrouded(block, `<li>`, `</li>`)
				blocks = blocks[1:]
			}
			fmt.Fprintf(g.f, `</ul>`)
		} else {
			g.genBlock(block)
			blocks = blocks[1:]
		}
	}
}

func (g *HTMLGenerator) genContent(parent *notionapi.Block) {
	g.genBlocks(parent.Content)
}