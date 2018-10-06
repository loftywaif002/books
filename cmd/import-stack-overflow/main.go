package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/essentialbooks/books/pkg/common"
	"github.com/essentialbooks/books/pkg/kvstore"
	"github.com/essentialbooks/books/pkg/stackoverflow"
	"github.com/gomarkdown/markdown"
	"github.com/kjk/u"
)

type DocTag = stackoverflow.DocTag
type Topic = stackoverflow.Topic
type Example = stackoverflow.Example
type TopicHistory = stackoverflow.TopicHistory
type Contributor = stackoverflow.Contributor

var (
	emptyExamplexs []*Example
	// if true, prints more information
	verbose = false

	booksToImport = common.BooksToProcess
)

func mdToHTML(d []byte) []byte {
	return markdown.ToHTML(d, nil, nil)
}

func getTopicsByDocID(docID int) map[int]bool {
	res := make(map[int]bool)
	topics := loadTopicsMust()
	for _, topic := range topics {
		if topic.DocTagId == docID {
			res[topic.Id] = true
		}
	}
	return res
}

func isEmptyString(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) == 0
}

func calcExampleCount(docTag *DocTag) {
	docID := docTag.Id
	topics := getTopicsByDocID(docID)
	n := 0
	examples := loadExamplesMust()
	for _, ex := range examples {
		if topics[ex.DocTopicId] {
			n++
		}
	}
	docTag.ExampleCount = n
}

func printDocTagsAndExit() {
	loadAll()
	docs := loadDocTagsMust()
	for i := range docs {
		docTag := &docs[i]
		calcExampleCount(docTag)
	}
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].ExampleCount < docs[j].ExampleCount
	})
	for _, dc := range docs {
		fmt.Printf(`{ "%s", "", false, %d, %d },%s`, dc.Title, dc.ExampleCount, dc.TopicCount, "\n")
	}
	os.Exit(0)
}

var (
	docTagsCached        []DocTag
	topicsCached         []Topic
	topicHistoriesCached []TopicHistory
	contributorsCached   []*Contributor
	examplesCached       []*Example
)

func loadDocTagsMust() []DocTag {
	if docTagsCached == nil {
		var err error
		path := path.Join("stack-overflow-docs-dump", "doctags.json.gz")
		docTagsCached, err = stackoverflow.LoadDocTags(path)
		u.PanicIfErr(err)
	}
	return docTagsCached
}

func loadTopicsMust() []Topic {
	if topicsCached == nil {
		var err error
		path := path.Join("stack-overflow-docs-dump", "topics.json.gz")
		topicsCached, err = stackoverflow.LoadTopics(path)
		u.PanicIfErr(err)
	}
	return topicsCached
}

func loadTopicHistoriesMust() []TopicHistory {
	if topicHistoriesCached == nil {
		var err error
		path := path.Join("stack-overflow-docs-dump", "topichistories.json.gz")
		topicHistoriesCached, err = stackoverflow.LoadTopicHistories(path)
		u.PanicIfErr(err)
	}
	return topicHistoriesCached
}

func loadContributorsMust() []*Contributor {
	if contributorsCached == nil {
		var err error
		path := path.Join("stack-overflow-docs-dump", "contributors.json.gz")
		contributorsCached, err = stackoverflow.LoadContibutors(path)
		u.PanicIfErr(err)
	}
	return contributorsCached
}

func loadExamplesMust() []*Example {
	if examplesCached == nil {
		var err error
		path := path.Join("stack-overflow-docs-dump", "examples.json.gz")
		examplesCached, err = stackoverflow.LoadExamples(path)
		u.PanicIfErr(err)
	}
	return examplesCached
}

func findDocTagByTitleMust(docTags []DocTag, title string) DocTag {
	for _, dc := range docTags {
		if dc.Title == title {
			return dc
		}
	}
	log.Fatalf("Didn't find DocTag with title '%s'\n", title)
	return DocTag{}
}

func loadAll() {
	timeStart := time.Now()
	fmt.Printf("Loading Stack Overflow data...")
	loadDocTagsMust()
	loadTopicsMust()
	loadExamplesMust()
	loadTopicHistoriesMust()
	loadContributorsMust()
	fmt.Printf(" took %s\n", time.Since(timeStart))
}

func getTopicsByDocTagID(docTagID int) []*Topic {
	gTopics := loadTopicsMust()
	var res []*Topic
	for i, topic := range gTopics {
		if topic.DocTagId == docTagID {
			res = append(res, &gTopics[i])
		}
	}
	return res
}

func getExampleByID(id int) *Example {
	gExamples := loadExamplesMust()
	for i, e := range gExamples {
		if e.Id == id {
			return gExamples[i]
		}
	}
	return nil
}

func getExamplesForTopic(docTagID int, docTopicID int) []*Example {
	gTopicHistories := loadTopicHistoriesMust()
	var res []*Example
	seenIds := make(map[int]bool)
	for _, th := range gTopicHistories {
		if th.DocTagId == docTagID && th.DocTopicId == docTopicID {
			id := th.DocExampleId
			if seenIds[id] {
				continue
			}
			seenIds[id] = true
			ex := getExampleByID(id)
			if ex == nil {
				//fmt.Printf("Didn't find example, docTagID: %d, docTopicID: %d\n", docTagID, docTopicID)
			} else {
				res = append(res, ex)
			}
		}
	}
	return res
}

func sortExamples(a []*Example) {
	sort.Slice(a, func(i, j int) bool {
		if a[i].IsPinned {
			return true
		}
		if a[j].IsPinned {
			return false
		}
		return a[i].Score > a[j].Score
	})
}

// sometime json representation of versions is empty array, we want to skip those
func shortenVersion(s string) string {
	if s == "[]" {
		return ""
	}
	return s
}

func writeIndexTxtMust(path string, topic *Topic) {
	s := kvstore.Serialize("Title", topic.Title)
	s += kvstore.Serialize("Id", strconv.Itoa(topic.Id))
	versions := shortenVersion(topic.VersionsJson)
	s += kvstore.SerializeLong("Versions", versions)
	if isEmptyString(versions) {
		s += kvstore.SerializeLong("VersionsHtml", topic.HelloWorldVersionsHtml)
	}

	s += kvstore.SerializeLong("Introduction", topic.IntroductionMarkdown)
	if isEmptyString(topic.IntroductionMarkdown) {
		s += kvstore.SerializeLong("IntroductionHtml", topic.IntroductionHtml)
	}

	s += kvstore.SerializeLong("Syntax", topic.SyntaxMarkdown)
	if isEmptyString(topic.SyntaxMarkdown) {
		s += kvstore.SerializeLong("SyntaxHtml", topic.SyntaxHtml)
	}

	s += kvstore.SerializeLong("Parameters", topic.ParametersMarkdown)
	if isEmptyString(topic.ParametersMarkdown) {
		s += kvstore.SerializeLong("ParametersHtml", topic.ParametersHtml)
	}

	s += kvstore.SerializeLong("Remarks", topic.RemarksMarkdown)
	if isEmptyString(topic.RemarksMarkdown) {
		s += kvstore.SerializeLong("RemarksHtml", topic.RemarksHtml)
	}

	createDirForFileMust(path)
	err := ioutil.WriteFile(path, []byte(s), 0644)
	u.PanicIfErr(err)
	if verbose {
		fmt.Printf("Wrote %s, %d bytes\n", path, len(s))
	}
}

func writeArticleMust(path string, example *Example) {
	s := kvstore.Serialize("Title", example.Title)
	s += kvstore.Serialize("Id", strconv.Itoa(example.Id))
	s += kvstore.Serialize("Score", strconv.Itoa(example.Score))
	s += kvstore.SerializeLong("Body", example.BodyMarkdown)
	if isEmptyString(example.BodyMarkdown) {
		s += kvstore.SerializeLong("BodyHtml", example.BodyHtml)
	}

	createDirForFileMust(path)
	err := ioutil.WriteFile(path, []byte(s), 0644)
	u.PanicIfErr(err)
	if verbose {
		fmt.Printf("Wrote %s, %d bytes\n", path, len(s))
	}
}

func printEmptyExamples() {
	for _, ex := range emptyExamplexs {
		fmt.Printf("empty example: %s, len(BodyHtml): %d\n", ex.Title, len(ex.BodyHtml))
	}
}

func getContributors(docID int) []int {
	gContributors := loadContributorsMust()
	topics := getTopicsByDocID(docID)
	contributors := make(map[int]bool)
	for _, c := range gContributors {
		topicID := c.DocTopicId
		if _, ok := topics[topicID]; ok {
			contributors[c.UserId] = true
		}
	}
	var res []int
	for id := range contributors {
		res = append(res, id)
	}
	return res
}

func genContributors(bookDstDir string, docID int) {
	contributors := getContributors(docID)
	var a []string
	for _, id := range contributors {
		a = append(a, strconv.Itoa(id))
	}
	s := strings.Join(a, "\n")
	path := filepath.Join(bookDstDir, "so_contributors.txt")
	createDirForFileMust(path)
	err := ioutil.WriteFile(path, []byte(s), 0644)
	u.PanicIfErr(err)
	//fmt.Printf("Wrote %s\n", path)
}

func importBook(docTag *DocTag, bookName string) {
	timeStart := time.Now()

	bookNameSafe := common.MakeURLSafe(bookName)
	bookTopDir := filepath.Join("books", bookNameSafe)
	if pathExists(bookTopDir) {
		fmt.Printf("Book '%s' has already been imported.\nTo re-import, delete directory '%s'\n", bookName, bookTopDir)
		os.Exit(1)
	}

	fmt.Printf("Importing a book %s\n", bookName)
	loadAll()

	//fmt.Printf("%s: docID: %d\n", title, docTag.Id)
	topics := getTopicsByDocTagID(docTag.Id)
	nChapters := len(topics)
	nArticles := 0
	chapter := 10
	for _, t := range topics {
		examples := getExamplesForTopic(docTag.Id, t.Id)
		sortExamples(examples)

		dirChapter := fmt.Sprintf("%04d-%s", chapter, common.MakeURLSafe(t.Title))
		dirPath := filepath.Join(bookTopDir, dirChapter)
		chapterIndexPath := filepath.Join(dirPath, "000-index.md")
		writeIndexTxtMust(chapterIndexPath, t)
		//fmt.Printf("%s\n", dirChapter)
		chapter += 10
		//fmt.Printf("%s, %d examples (%d), %s\n", t.Title, t.ExampleCount, len(examples), fileName)

		articleNo := 10
		for _, ex := range examples {
			if isEmptyString(ex.BodyMarkdown) && isEmptyString(ex.BodyHtml) {
				emptyExamplexs = append(emptyExamplexs, ex)
				continue
			}
			fileName := fmt.Sprintf("%03d-%s.md", articleNo, common.MakeURLSafe(ex.Title))
			path := filepath.Join(dirPath, fileName)
			writeArticleMust(path, ex)
			//fmt.Printf("  %s %s '%s'\n", ex.Title, pinnedStr, fileName)
			//fmt.Printf("  %03d-%s\n", articleNo, fileName)
			//fmt.Printf("  %s\n", fileName)
			articleNo += 10
		}
		nArticles += len(examples)
	}
	genContributors(bookTopDir, docTag.Id)

	fmt.Printf("Imported %s (%d chapters, %d articles) in %s\n", bookName, nChapters, nArticles, time.Since(timeStart))
}

func dumpMetaAndExit() {
	loadAll()
	os.Exit(0)
}

func getImportedBooks() []string {
	books, err := common.GetDirs("books")
	u.PanicIfErr(err)
	return books
}

func getAllBooks() []string {
	gDocTags := loadDocTagsMust()
	var books []string
	for _, doc := range gDocTags {
		book := doc.Title
		books = append(books, book)
	}
	return books
}

func printAllBookNames() {
	all := getAllBooks()
	s := strings.Join(all, ", ")
	fmt.Printf("All books: %s\n", s)
}

func printUsageAndExit() {
	fmt.Printf("Usage: import-stack-overflow book-to-import\n")
	imported := getImportedBooks()
	if len(imported) > 0 {
		s := strings.Join(imported, ", ")
		fmt.Printf("Already imported: %s\n", s)
	}
	printAllBookNames()
	os.Exit(1)
}

// some custom fixups for stack overflow book name => our book name
var bookNameFixups = [][]string{
	{"Intel x86 Assembly Language & Microarchitecture", "Intel x86 assembly"},
	{"tensorflow", "TensorFlow"},
	{"react-native", "React Native"},
	{"postgresql", "PostgreSQL"},
	{"batch-file", "Batch file"},
	{"excel-vba", "Excel VBA"},
	{"html5-canvas", "HTML Canvas"},
	{"algorithm", "Algorithms"},
	{"meteor", "Meteor"},
}

// convert a stack overflow book name to our book name
func fixupBookName(soName string) string {
	// "Ruby Language" => "Ruby" etc.
	if strings.HasSuffix(soName, "Language") {
		s := strings.TrimSuffix(soName, "Language")
		return strings.TrimSpace(s)
	}
	// manual overrides
	for _, fixup := range bookNameFixups {
		if soName == fixup[0] {
			return fixup[1]
		}
	}
	return soName
}

func findBookByName(bookName string) *DocTag {
	nameNoCase := strings.ToLower(bookName)
	gDocTags := loadDocTagsMust()
	for i, doc := range gDocTags {
		titleNoCase := strings.ToLower(doc.Title)
		if nameNoCase == titleNoCase {
			return &gDocTags[i]
		}
	}
	return nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func main() {
	// for ad-hoc operations uncomment one of those
	// genContributorsAndExit()
	// dumpMetaAndExit()
	// printDocTagsAndExit()

	args := os.Args[1:]
	if len(args) != 1 {
		printUsageAndExit()
	}
	timeStart := time.Now()
	fmt.Printf("Trying to import book %s\n", args[0])

	bookName := args[0]
	doc := findBookByName(bookName)
	if doc == nil {
		printAllBookNames()
		fmt.Printf("\nDidn't find a book '%s'.\nSee above for list of available books\n", bookName)
		os.Exit(1)
	}

	importBook(doc, bookName)

	fmt.Printf("Took %s\n", time.Since(timeStart))
	//printEmptyExamples()
}
