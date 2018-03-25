package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/essentialbooks/books/pkg/common"
	"github.com/essentialbooks/books/pkg/kvstore"
	"github.com/kjk/u"
)

var (
	defTitle = "No Title"

	bookDirToName = map[string]string{
		"go": "Go",
		"Go": "Go",
	}
)

func dumpKV(doc kvstore.Doc) {
	for _, kv := range doc {
		fmt.Printf("K: %s\nV: %s\n", kv.Key, common.ShortenString(kv.Value))
	}
}

func parseArticle(path string) (*Article, error) {
	kvdoc, err := parseKVFileWithIncludes(path)
	if err != nil {
		fmt.Printf("Error parsing KV file: '%s'\n", path)
		maybePanicIfErr(err)
		return nil, err
	}

	doc := &MarkdownFile{
		Path: path,
	}
	article := &Article{
		MarkdownFile: doc,
	}
	article.ID, err = kvdoc.Get("Id")
	if err != nil {
		return nil, fmt.Errorf("parseArticle('%s'), err: '%s'", path, err)
	}
	if strings.Contains(article.ID, " ") {
		return nil, fmt.Errorf("parseArticle('%s'), res.ID = '%s' has space in it", path, article.ID)
	}

	article.Title = kvdoc.GetSilent("Title", defTitle)
	if article.Title == defTitle {
		fmt.Printf("parseArticle: no title for %s\n", path)
	}
	titleSafe := common.MakeURLSafe(article.Title)

	// handle search synonyms
	synonyms := kvdoc.GetSilent("Search", "")
	synonyms = strings.TrimSpace(synonyms)
	if len(synonyms) > 0 {
		parts := strings.Split(synonyms, ",")
		for _, synonym := range parts {
			synonym = strings.TrimSpace(synonym)
			if len(synonym) > 0 {
				article.SearchSynonyms = append(article.SearchSynonyms, synonym)
			}
		}
	}

	article.FileNameBase = fmt.Sprintf("%s-%s", article.ID, titleSafe)
	article.BodyMarkdown, err = kvdoc.Get("Body")
	if err == nil {
		return article, nil
	}
	s, err := kvdoc.Get("BodyHtml")
	article.BodyHTML = template.HTML(s)
	if err != nil {
		dumpKV(kvdoc)
		return nil, fmt.Errorf("parseArticle('%s'), err: '%s'", path, err)
	}
	return article, nil
}

func buildArticleSiblings(articles []*Article) {
	// build a template
	var siblings []Article
	for i, article := range articles {
		sibling := *article // making a copy, we can't touch the original
		sibling.No = i + 1
		siblings = append(siblings, sibling)
	}
	// for each article, copy a template and set IsCurrent
	for i, article := range articles {
		copy := append([]Article(nil), siblings...)
		copy[i].IsCurrent = true
		article.Siblings = copy
	}
}

// Parses @file ${fileName} directives and replaces them
// with the content of the file
func processFileIncludes(path string) ([]string, error) {
	fc, err := loadFileCached(path)
	if err != nil {
		return nil, err
	}
	lines := fc.Lines
	nLines := len(lines)
	res := make([]string, 0, nLines)
	for _, line := range lines {
		if !strings.HasPrefix(line, "@file") {
			res = append(res, line)
			continue
		}

		//fmt.Printf("processFileIncludes('%s'\n", path)
		lines2, err := extractCodeSnippetsAsMarkdownLines(filepath.Dir(path), line)
		if err != nil {
			fmt.Printf("processFileIncludes: error '%s'\n", err)
			return nil, err
		}
		res = append(res, lines2...)
	}
	return res, nil
}

func parseKVFileWithIncludes(path string) (kvstore.Doc, error) {
	lines, err := processFileIncludes(path)
	if err == nil {
		return kvstore.ParseKVLines(lines)
	}
	// if processFileIncludes fails we retry without file includes
	return kvstore.ParseKVFile(path)
}

func parseChapter(chapter *Chapter) error {
	dir := filepath.Join(chapter.Book.sourceDir, chapter.ChapterDir)
	path := filepath.Join(dir, "000-index.md")
	chapter.Path = path
	doc, err := parseKVFileWithIncludes(path)
	if err != nil {
		fmt.Printf("Error parsing KV file: '%s'\n", path)
		maybePanicIfErr(err)
	}

	chapter.indexDoc = doc
	chapter.Title, err = doc.Get("Title")
	if err != nil {
		return fmt.Errorf("parseChapter('%s'), missing Title, err: '%s'", path, err)
	}
	chapter.ID, err = doc.Get("Id")
	if err != nil {
		return fmt.Errorf("parseChapter('%s'), missing Id, err: '%s'", path, err)
	}

	if strings.Contains(chapter.ID, " ") {
		return fmt.Errorf("parseChapter('%s'), chapter.ID = '%s' has space in it", path, chapter.ID)
	}

	titleSafe := common.MakeURLSafe(chapter.Title)
	chapter.FileNameBase = fmt.Sprintf("%s-%s", chapter.ID, titleSafe)
	fileInfos, err := ioutil.ReadDir(dir)
	var articles []*Article
	for _, fi := range fileInfos {
		if fi.IsDir() || !fi.Mode().IsRegular() {
			continue
		}
		name := fi.Name()
		ext := strings.ToLower(filepath.Ext(name))

		// remember images to be copied in gen book phase
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
			path = filepath.Join(dir, name)
			chapter.images = append(chapter.images, path)
			continue
		}

		if ext != ".md" {
			continue
		}

		// some files are not meant to be processed here
		switch strings.ToLower(name) {
		case "000-index.md":
			continue
		}
		path = filepath.Join(dir, name)
		article, err := parseArticle(path)
		if err != nil {
			return err
		}
		article.Chapter = chapter
		article.No = len(articles) + 1
		articles = append(articles, article)
	}
	buildArticleSiblings(articles)
	chapter.Articles = articles
	return nil
}

func soContributorURL(userID int, userName string) string {
	return fmt.Sprintf("https://stackoverflow.com/users/%d/%s", userID, userName)
}

func loadSoContributorsMust(book *Book, path string) {
	fc, err := loadFileCached(path)
	u.PanicIfErr(err)
	lines := fc.Lines
	var contributors []SoContributor
	for _, line := range lines {
		id, err := strconv.Atoi(line)
		u.PanicIfErr(err)
		name := soUserIDToNameMap[id]
		u.PanicIf(name == "", "no SO contributor for id %d", id)
		if name == "user_deleted" {
			continue
		}
		nameUnescaped, err := url.PathUnescape(name)
		u.PanicIfErr(err)
		c := SoContributor{
			ID:      id,
			URLPart: name,
			Name:    nameUnescaped,
		}
		contributors = append(contributors, c)
	}
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Name < contributors[j].Name
	})
	book.SoContributors = contributors
}

// TODO: add github contributors
func genContributorsMarkdown(contributors []SoContributor) string {
	if len(contributors) == 0 {
		return ""
	}
	lines := []string{
		"Contributors from [GitHub](https://github.com/essentialbooks/books/graphs/contributors)",
		"",
		"Contributors from Stack Overflow:",
	}
	for _, c := range contributors {
		s := fmt.Sprintf("* [%s](%s)", c.Name, soContributorURL(c.ID, c.Name))
		lines = append(lines, s)
	}
	return strings.Join(lines, "\n")
}

func genContributorsChapter(book *Book) *Chapter {
	md := genContributorsMarkdown(book.SoContributors)
	var kvdoc kvstore.Doc
	kv := kvstore.KeyValue{
		Key:   "Body",
		Value: md,
	}
	kvdoc = append(kvdoc, kv)
	doc := &MarkdownFile{
		No:           999,
		Title:        "Contributors",
		FileNameBase: "contributors",
	}
	ch := &Chapter{
		MarkdownFile: doc,
		Book:         book,
		indexDoc:     kvdoc,
	}
	return ch
}

// make sure chapter/article ids within the book are unique,
// so that we can generate stable urls.
// also build a list of chapter/article urls
func ensureUniqueIds(book *Book) {
	var urls []string
	chapterIds := make(map[string]*Chapter)
	articleIds := make(map[string]*Article)
	for _, c := range book.Chapters {
		if chap, ok := chapterIds[c.ID]; ok {
			fmt.Printf("Duplicate chapter id '%s' in:\n", c.ID)
			fmt.Printf("Chapter '%s', file: '%s'\n", c.Title, c.Path)
			fmt.Printf("Chapter '%s', file: '%s'\n", chap.Title, chap.Path)
			os.Exit(1)
		}
		chapterIds[c.ID] = c
		urls = append(urls, c.FileNameBase)
		for _, a := range c.Articles {
			if a2, ok := articleIds[a.ID]; ok {
				err := fmt.Errorf("Duplicate article id: '%s', in: %s and %s", a.ID, a.Path, a2.Path)
				maybePanicIfErr(err)
			} else {
				articleIds[a.ID] = a
				urls = append(urls, a.FileNameBase)
			}
		}
	}
	book.knownUrls = urls
}

func parseBook(bookDir string) (*Book, error) {
	timeStart := time.Now()
	bookName := bookDir
	bookName, ok := bookDirToName[bookDir]
	u.PanicIf(!ok, "no book name from dir '%s'", bookDir)
	fmt.Printf("Parsing book %s\n", bookName)
	bookNameSafe := common.MakeURLSafe(bookName)
	srcDir := filepath.Join("books", bookNameSafe)
	book := &Book{
		Title:        bookName,
		titleSafe:    bookNameSafe,
		TitleLong:    fmt.Sprintf("Essential %s", bookName),
		FileNameBase: bookNameSafe,
		sourceDir:    srcDir,
		destDir:      filepath.Join(destEssentialDir, bookNameSafe),
	}

	fileInfos, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return nil, err
	}

	nProcs := getAlmostMaxProcs()

	sem := make(chan bool, nProcs)
	var wg sync.WaitGroup
	var chapters []*Chapter
	var err2 error

	for _, fi := range fileInfos {
		if fi.IsDir() {
			ch := &Chapter{
				Book:       book,
				ChapterDir: fi.Name(),
			}
			chapters = append(chapters, ch)
			sem <- true
			wg.Add(1)
			go func(chap *Chapter) {
				err = parseChapter(chap)
				if err != nil {
					// not thread safe but whatever
					err2 = err
				}
				<-sem
				wg.Done()
			}(ch)
			continue
		}

		name := strings.ToLower(fi.Name())
		// some files should be ignored
		if name == "toc.txt" {
			continue
		}
		if name == "so_contributors.txt" {
			path := filepath.Join(srcDir, fi.Name())
			loadSoContributorsMust(book, path)
			continue
		}
		return nil, fmt.Errorf("Unexpected file at top-level: '%s'", fi.Name())
	}
	wg.Wait()

	ch := genContributorsChapter(book)
	chapters = append(chapters, ch)

	for i, ch := range chapters {
		ch.No = i + 1
	}
	book.Chapters = chapters

	ensureUniqueIds(book)

	fmt.Printf("Book '%s' %d chapters, %d articles, finished parsing in %s\n", bookName, len(chapters), book.ArticlesCount(), time.Since(timeStart))
	return book, err2
}
