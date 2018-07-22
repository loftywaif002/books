package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kjk/notionapi"
)

var (
	useCacheForNotion = true
	// if true, we'll log
	logNotionRequests = true

	notionIDToPage = map[string]*notionapi.PageInfo{}

	cacheDir     = "notion_cache"
	notionLogDir = "log"

	// "https://www.notion.so/kjkpublic/Essential-Go-2cab1ed2b7a44584b56b0d3ca9b80185"
	notionGoStartPage = "2cab1ed2b7a44584b56b0d3ca9b80185"
)

// convert 2131b10c-ebf6-4938-a127-7089ff02dbe4 to 2131b10cebf64938a1277089ff02dbe4
func normalizeID(s string) string {
	return strings.Replace(s, "-", "", -1)
}

func openLogFileForPageID(pageID string) (io.WriteCloser, error) {
	if !logNotionRequests {
		return nil, nil
	}

	name := fmt.Sprintf("%s.go.log.txt", pageID)
	path := filepath.Join(notionLogDir, name)
	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("os.Create('%s') failed with %s\n", path, err)
		return nil, err
	}
	notionapi.Logger = f
	return f, nil
}

func findSubPageIDs(blocks []*notionapi.Block) []string {
	pageIDs := map[string]struct{}{}
	seen := map[string]struct{}{}
	toVisit := blocks
	for len(toVisit) > 0 {
		block := toVisit[0]
		toVisit = toVisit[1:]
		id := normalizeID(block.ID)
		if block.Type == notionapi.BlockPage {
			pageIDs[id] = struct{}{}
			seen[id] = struct{}{}
		}
		for _, b := range block.Content {
			if b == nil {
				continue
			}
			id := normalizeID(block.ID)
			if _, ok := seen[id]; ok {
				continue
			}
			toVisit = append(toVisit, b)
		}
	}
	res := []string{}
	for id := range pageIDs {
		res = append(res, id)
	}
	sort.Strings(res)
	return res
}

func loadPageFromCache(pageID string) *notionapi.PageInfo {
	if !useCacheForNotion {
		return nil
	}

	cachedPath := filepath.Join(cacheDir, pageID+".json")
	d, err := ioutil.ReadFile(cachedPath)
	if err != nil {
		return nil
	}

	var pageInfo notionapi.PageInfo
	err = json.Unmarshal(d, &pageInfo)
	panicIfErr(err)
	fmt.Printf("Got %s from cache (%s)\n", pageID, pageInfo.Page.Title)
	return &pageInfo
}

func downloadAndCachePage(pageID string) (*notionapi.PageInfo, error) {
	//fmt.Printf("downloading page with id %s\n", pageID)
	lf, _ := openLogFileForPageID(pageID)
	if lf != nil {
		defer lf.Close()
	}
	cachedPath := filepath.Join(cacheDir, pageID+".json")
	res, err := notionapi.GetPageInfo(pageID)
	if err != nil {
		return nil, err
	}
	d, err := json.MarshalIndent(res, "", "  ")
	if err == nil {
		err = ioutil.WriteFile(cachedPath, d, 0644)
		panicIfErr(err)
	} else {
		// not a fatal error, just a warning
		fmt.Printf("json.Marshal() on pageID '%s' failed with %s\n", pageID, err)
	}
	return res, nil
}

func loadNotionPages(indexPageID string) []*notionapi.PageInfo {
	var res []*notionapi.PageInfo

	toVisit := []string{indexPageID}

	for len(toVisit) > 0 {
		pageID := normalizeID(toVisit[0])
		toVisit = toVisit[1:]

		if _, ok := notionIDToPage[pageID]; ok {
			continue
		}

		var err error
		page := loadPageFromCache(pageID)
		if page == nil {
			page, err = downloadAndCachePage(pageID)
			panicIfErr(err)
			fmt.Printf("Downloaded %s %s\n", pageID, page.Page.Title)
		}

		notionIDToPage[pageID] = page
		res = append(res, page)

		subPages := findSubPageIDs(page.Page.Content)
		toVisit = append(toVisit, subPages...)
	}

	return res
}

func loadAllPages() []*notionapi.PageInfo {
	loadNotionPages(notionGoStartPage)
	n := len(notionIDToPage)
	res := make([]*notionapi.PageInfo, 0, n)
	for _, page := range notionIDToPage {
		res = append(res, page)
	}
	return res
}

func rmFile(path string) {
	err := os.Remove(path)
	if err != nil {
		fmt.Printf("os.Remove(%s) failed with %s\n", path, err)
	}
}

func rmCached(pageID string) {
	id := normalizeID(pageID)
	rmFile(filepath.Join(notionLogDir, id+".go.log.txt"))
	rmFile(filepath.Join(cacheDir, id+".json"))
}

func createNotionDirs() {
	if logNotionRequests {
		err := os.MkdirAll(notionLogDir, 0755)
		panicIfErr(err)
	}
	{
		err := os.MkdirAll(cacheDir, 0755)
		panicIfErr(err)
	}
}

func maybeRemoveNotionCache() {
	if useCacheForNotion {
		return
	}
	err := os.RemoveAll(cacheDir)
	panicIfErr(err)
	err = os.RemoveAll(notionLogDir)
	panicIfErr(err)
}

// this re-downloads pages from Notion by deleting cache locally
func notionRedownload() {
	//notionapi.DebugLog = true
	maybeRemoveNotionCache()
	createNotionDirs()

	pages := loadAllPages()
	fmt.Printf("Loaded %d pages\n", len(pages))
}
