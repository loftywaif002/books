package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/essentialbooks/books/pkg/common"
	"github.com/kjk/u"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

var (
	flgPreview            bool
	flgUpdateGoPlayground bool
	flgUpdateOutput       bool
	flgRecreateOutput     bool
	flgForce              bool
	flgUpdateGoDeps       bool
	flgGenID              bool

)

func parseFlags() {
	flag.BoolVar(&flgPreview, "preview", false, "if true will start watching for file changes and re-build everything")
	flag.BoolVar(&flgUpdateGoPlayground, "update-go-playground", false, "if true will upgrade links to go playground")
	flag.BoolVar(&flgUpdateOutput, "update-output", false, "if true, will update ouput files in cached_output")
	flag.BoolVar(&flgRecreateOutput, "recreate-output", false, "if true, recreates ouput files in cached_output")
	flag.BoolVar(&flgUpdateGoDeps, "update-go-deps", false, "if true, updates go libraries references in go snippets")
	flag.BoolVar(&flgGenID, "gen-id", false, "if true, generate unique id")
	flag.Parse()

	if flgAnalytics != "" {
		s := fmt.Sprintf(googleAnalyticsTmpl, flgAnalytics, flgAnalytics)
		googleAnalytics = template.HTML(s)
	}
}

func dirFromBook(book *common.Book) string {
	return common.MakeURLSafe(book.NewName())
}

func isBookImported(bookDirs []string, book *common.Book) bool {
	dir := dirFromBook(book)
	for _, dir2 := range bookDirs {
		if dir == dir2 {
			return true
		}
	}
	return false
}

func getBooksToImport(bookDirs []string) []*common.Book {
	var res []*common.Book
	for _, book := range common.BooksToProcess {
		if isBookImported(bookDirs, book) {
			res = append(res, book)
		}
	}
	return res
}

// TODO: probably more
func getDefaultLangForBook(bookName string) string {
	s := strings.ToLower(bookName)
	switch s {
	case "go":
		return "go"
	case "android":
		return "java"
	case "ios":
		return "ObjectiveC"
	case "microsoft sql server":
		return "sql"
	case "node.js":
		return "javascript"
	case "mysql":
		return "sql"
	case ".net framework":
		return "c#"
	}
	return s
}

func getBookDirs() []string {
	dirs, err := common.GetDirs("books")
	u.PanicIfErr(err)
	return dirs
}

func shouldCopyImage(path string) bool {
	return !strings.Contains(path, "@2x")
}

func copyFilesRecur(dstDir, srcDir string, shouldCopyFunc func(path string) bool) {
	createDirMust(dstDir)
	fileInfos, err := ioutil.ReadDir(srcDir)
	u.PanicIfErr(err)
	for _, fi := range fileInfos {
		name := fi.Name()
		if fi.IsDir() {
			dst := filepath.Join(dstDir, name)
			src := filepath.Join(srcDir, name)
			copyFilesRecur(dst, src, shouldCopyFunc)
			continue
		}

		src := filepath.Join(srcDir, name)
		dst := filepath.Join(dstDir, name)
		shouldCopy := true
		if shouldCopyFunc != nil {
			shouldCopy = shouldCopyFunc(src)
		}
		if !shouldCopy {
			continue
		}
		if pathExists(dst) {
			continue
		}
		copyFileMust(dst, src)
	}
}

func copyCoversMust() {
	copyFilesRecur(filepath.Join("www", "covers"), "covers", shouldCopyImage)
}

func getAlmostMaxProcs() int {
	// leave some juice for other programs
	nProcs := runtime.NumCPU() - 2
	if nProcs < 1 {
		return 1
	}
	return nProcs
}

// copy from tmpl to www, optimize if possible, add
// sha1 of the content as part of the name
func copyToWwwAsSha1MaybeMust(srcName string) {
	var dstPtr *string
	minifyType := ""
	switch srcName {
	case "main.css":
		dstPtr = &pathMainCSS
		minifyType = "text/css"
	case "app.js":
		dstPtr = &pathAppJS
		minifyType = "text/javascript"
	case "favicon.ico":
		dstPtr = &pathFaviconICO
	default:
		u.PanicIf(true, "unknown srcName '%s'", srcName)
	}
	src := filepath.Join("tmpl", srcName)
	d, err := ioutil.ReadFile(src)
	u.PanicIfErr(err)

	if doMinify && minifyType != "" {
		d2, err := minifier.Bytes(minifyType, d)
		maybePanicIfErr(err)
		if err == nil {
			fmt.Printf("Compressed %s from %d => %d (saved %d)\n", srcName, len(d), len(d2), len(d)-len(d2))
			d = d2
		}
	}

	sha1Hex := u.Sha1HexOfBytes(d)
	name := nameToSha1Name(srcName, sha1Hex)
	dst := filepath.Join("www", "s", name)
	err = ioutil.WriteFile(dst, d, 0644)
	u.PanicIfErr(err)
	*dstPtr = filepath.ToSlash(dst[len("www"):])
	fmt.Printf("Copied %s => %s\n", src, dst)
}

func genSelectedBooks(bookDirs []string) {
	fmt.Printf("genSelectedBooks: %+v\n", bookDirs)
	timeStart := time.Now()

	var books []*Book
	for _, bookName := range bookDirs {
		book, err := parseBook(bookName)
		maybePanicIfErr(err)
		if err != nil {
			continue
		}
		book.sem = make(chan bool, getAlmostMaxProcs())
		books = append(books, book)
	}
	fmt.Printf("Parsed books in %s\n", time.Since(timeStart))

	copyToWwwAsSha1MaybeMust("main.css")
	copyToWwwAsSha1MaybeMust("app.js")
	copyToWwwAsSha1MaybeMust("favicon.ico")
	genIndex(books)
	genIndexGrid(books)
	gen404TopLevel()
	genAbout()
	genFeedback()

	for _, book := range books {
		genBook(book)
	}
	fmt.Printf("Used %d procs, finished generating all books in %s\n", getAlmostMaxProcs(), time.Since(timeStart))
}

func genAllBooks(udpateOutputCache bool) {
	timeStart := time.Now()
	clearSitemapURLS()
	copyCoversMust()

	nProcs := getAlmostMaxProcs()

	var books []*Book
	for _, bookName := range allBookDirs {
		book, err := parseBook(bookName)
		maybePanicIfErr(err)
		if err != nil {
			continue
		}
		book.sem = make(chan bool, nProcs)
		books = append(books, book)
	}
	fmt.Printf("Parsed books in %s\n", time.Since(timeStart))

	copyToWwwAsSha1MaybeMust("main.css")
	copyToWwwAsSha1MaybeMust("app.js")
	copyToWwwAsSha1MaybeMust("favicon.ico")
	genIndex(books)
	gen404TopLevel()
	genIndexGrid(books)
	genAbout()
	genFeedback()

	for _, book := range books {
		genBook(book)
	}
	writeSitemap()
	if udpateOutputCache {
		saveCachedOutputFiles()
	}
	fmt.Printf("Used %d procs, finished generating all books in %s\n", nProcs, time.Since(timeStart))
}

func loadSOUserMappingsMust() {
	path := filepath.Join("stack-overflow-docs-dump", "users.json.gz")
	err := common.JSONDecodeGzipped(path, &soUserIDToNameMap)
	u.PanicIfErr(err)
}

func genNetlifyHeaders() {
	path := filepath.Join("www", "_headers")
	err := ioutil.WriteFile(path, []byte(netlifyHeaders), 0644)
	u.PanicIfErr(err)
}

func genNetlifyRedirects() {
	// TODO: should be generated from []*Book list
	s := `
/essential/go/* /essential/go/404.html 404
`
	path := filepath.Join("www", "_redirects")
	err := ioutil.WriteFile(path, []byte(s), 0644)
	u.PanicIfErr(err)
}

var intIDS = make(map[int]bool)

func rememberID(id string) {
	intID, err := strconv.Atoi(id)
	u.PanicIfErr(err, "'%s' id is not an int", id)
	if intIDS[intID] {
		fmt.Printf("duplicate id: %d\n", intID)
		os.Exit(1)
	}
	intIDS[intID] = true
}

func genID() {
	for _, bookName := range allBookDirs {
		book, err := parseBook(bookName)
		u.PanicIfErr(err)
		for _, chapter := range book.Chapters {
			if chapter.FileNameBase == "contributors" {
				continue
			}
			rememberID(chapter.ID)
			for _, article := range chapter.Articles {
				rememberID(article.ID)
			}
		}
	}
	var idArr []int
	for id := range intIDS {
		idArr = append(idArr, id)
	}
	sort.Ints(idArr)
	n := len(idArr)
	lastID := idArr[n-1]
	newID := lastID + 1
	//fmt.Printf("%v\n", idArr)
	fmt.Printf("id: %d\n", newID)
}

func main() {

	parseFlags()

	if false {
		genTwitterImagesAndExit()
	}

	if false {
		testGetGoPlaygroundShareIDAndExit()
	}

	if flgUpdateGoPlayground {
		goBookDir := filepath.Join("books", "go")
		updateGoPlaygroundLinks(goBookDir)
		os.Exit(0)
	}

	minifier = minify.New()
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("text/html", html.Minify)
	minifier.AddFunc("text/javascript", js.Minify)
	// less aggresive minification because html validators
	// report this as html errors
	minifier.Add("text/html", &html.Minifier{
		KeepDocumentTags: true,
		KeepEndTags:      true,
	})
	doMinify = !flgPreview

	reloadCachedOutputFilesMust()

	booksToImport := getBooksToImport(getBookDirs())
	for _, bookInfo := range booksToImport {
		allBookDirs = append(allBookDirs, bookInfo.NewName())
	}
	loadSOUserMappingsMust()

	if flgGenID {
		genID()
		os.Exit(0)
	}

	os.RemoveAll("www")
	createDirMust(filepath.Join("www", "s"))
	genNetlifyHeaders()
	genNetlifyRedirects()

	if flgUpdateGoDeps {
		updateGoDeps()
	}

	cacheFilesInDir("books")

	if flgUpdateOutput {
		gitRemoveCachedOutputFiles()
	}

	clearErrors()
	genAllBooks(flgUpdateOutput)
	printAndClearErrors()
	if flgUpdateOutput {
		gitAddachedOutputFiles()
		return
	}

	if flgPreview {
		startPreview()
	}
}
