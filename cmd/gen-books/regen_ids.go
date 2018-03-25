package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/essentialbooks/books/pkg/kvstore"
	"github.com/kjk/u"
)

type modifiedDoc struct {
	doc  kvstore.Doc
	path string
}

// this is a one-time operation to regenerate ids to make them shorter
func regenIDSAndExit() {
	// maps old StackOverflow id to new ids
	idMap := make(map[string]string)
	var docs []*modifiedDoc
	booksToImport := getBooksToImport(getBookDirs())
	for _, bookInfo := range booksToImport {
		allBookDirs = append(allBookDirs, bookInfo.NewName())
	}
	loadSOUserMappingsMust()
	for _, bookName := range allBookDirs {
		book, err := parseBook(bookName)
		u.PanicIfErr(err)
		for _, chapter := range book.Chapters {
			if chapter.FileNameBase == "contributors" {
				continue
			}

			path := chapter.Path
			doc, err := kvstore.ParseKVFile(path)
			u.PanicIfErr(err)

			mdoc := &modifiedDoc{
				doc:  doc,
				path: path,
			}
			docs = append(docs, mdoc)

			for _, article := range chapter.Articles {
				path := article.Path
				doc, err := kvstore.ParseKVFile(path)
				u.PanicIfErr(err)

				mdoc := &modifiedDoc{
					doc:  doc,
					path: path,
				}
				docs = append(docs, mdoc)
			}
		}
	}

	// assign new ids
	for id, mdoc := range docs {
		doc := mdoc.doc
		soID, err := doc.Get("Id")
		u.PanicIfErr(err)
		doc = kvstore.ReplaceOrAppend(doc, "SOId", soID)
		newID := strconv.Itoa(id + 1)
		mdoc.doc = kvstore.ReplaceOrAppend(doc, "Id", newID)
		idMap[soID] = newID
	}

	// update links in the body
	var old []string
	for k := range idMap {
		old = append(old, k)
	}
	sort.Slice(old, func(i, j int) bool {
		s1 := old[i]
		i1, err1 := strconv.Atoi(s1)
		s2 := old[j]
		i2, err2 := strconv.Atoi(s2)

		// both are integers
		if err1 == nil && err2 == nil {
			return i1 < i2
		}
		// both are string
		if err1 != nil && err2 != nil {
			return s1 < s2
		}
		if err1 != nil {
			// first value is string so bigger than second value
			return false
		}
		return true
	})

	// fix links (${oldID}) => (${newID})
	for _, oldID := range old {
		newID := idMap[oldID]
		fmt.Printf("%s => %s\n", oldID, newID)
		oldIDLink := "(" + oldID + ")"
		newIDLink := "(" + newID + ")"

		for _, mdoc := range docs {
			doc := mdoc.doc
			body := doc.GetSilent("Body", "")
			newBody := strings.Replace(body, oldIDLink, newIDLink, -1)
			mdoc.doc = kvstore.ReplaceOrAppend(doc, "Body", newBody)
		}
	}

	for _, mdoc := range docs {
		err := saveDoc(mdoc.path, mdoc.doc)
		u.PanicIfErr(err)
	}
	os.Exit(0)
}

func saveDoc(path string, doc kvstore.Doc) error {
	s, err := kvstore.SerializeDoc(doc)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(s), 0644)
}
