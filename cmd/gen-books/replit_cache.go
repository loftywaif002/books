package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/kjk/siser"
)

// ReplitFile describes a single file in a multi-file replit
type ReplitFile struct {
	name string
	data string
}

// Replit describes a single replit. It's a set of files.
type Replit struct {
	// url of replit e.g. https://repl.it/@kjk1/inflect-examples
	url   string
	files []*ReplitFile
}

// SortFiles sorts files so that it's easy to compare them
func (r *Replit) SortFiles() {
	sort.Slice(r.files, func(i, j int) bool {
		return r.files[i].name < r.files[j].name
	})
}

// Equal returns true if replits have the same content. We can only compare
// replits with the same url.
// Both replits should be already sorted
func (r *Replit) Equal(r2 *Replit) bool {
	panicIf(r.url != r2.url, "comparing replits with different urls. '%s' != '%s'", r.url, r2.url)

	n := len(r.files)
	n2 := len(r2.files)
	if n != n2 {
		return false
	}

	for i := 0; i < n; i++ {
		if r.files[i].name != r2.files[i].name {
			return false
		}
		if r.files[i].data != r2.files[i].data {
			return false
		}
	}

	return true
}

// ReplitCache represents a cache for multiple replits. We have one cache
// per each book
type ReplitCache struct {
	path string
	f    *os.File
	// maps a replit url to Replit
	replits map[string]*Replit
}

// must have @url key which is url of the replit
// other keys are file names with value being their content
func recordToReplit(r *siser.Record) (*Replit, error) {
	res := &Replit{}
	n := len(r.Keys)
	for i := 0; i < n; i++ {
		k := r.Keys[i]
		v := r.Values[i]
		if k == "@url" {
			res.url = v
			continue
		}
		rf := &ReplitFile{
			name: k,
			data: v,
		}
		res.files = append(res.files, rf)
	}
	if res.url == "" {
		return nil, fmt.Errorf("siser.Record is missing @url key. Keys: %#v", r.Keys)
	}
	res.SortFiles()
	return res, nil
}

func replitToRecord(r *Replit) *siser.Record {
	var rec siser.Record
	rec.Append("@url", r.url)
	for _, rf := range r.files {
		rec.Append(rf.name, rf.data)
	}
	return &rec
}

// LoadReplitCache loads existing cache and positions it for appends
// it'll. Cache doesn't have to exist.
func LoadReplitCache(path string) (*ReplitCache, error) {
	res := &ReplitCache{
		path:    path,
		replits: map[string]*Replit{},
	}
	f, err := os.Open(path)
	if err == nil {
		r := siser.NewReader(f)
		for r.ReadNext() {
			_, rec := r.Record()
			replit, err := recordToReplit(rec)
			if err != nil {
				f.Close()
				return nil, err
			}
			res.replits[replit.url] = replit
		}
		err := r.Err()
		f.Close()
		if err != nil {
			return nil, err
		}
	}
	res.f, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Add adds a replit to a cache. Returns true if new replit or different than
// existing replit
func (c *ReplitCache) Add(r *Replit) (bool, error) {
	re := c.replits[r.url]
	// ignore if exactly the same replit already exists
	if re != nil {
		if re.Equal(r) {
			return false, nil
		}
	}
	rec := replitToRecord(r)
	d := rec.Marshal()
	_, err := c.f.Write(d)
	if err != nil {
		return false, err
	}
	c.replits[r.url] = r
	return true, nil
}

// Close closes a file used for storing this cache
func (c *ReplitCache) Close() error {
	if c.f != nil {
		err := c.f.Close()
		c.f = nil
		return err
	}
	c.replits = nil
	return nil
}

func httpGet(uri string) ([]byte, error) {
	hc := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := hc.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		d, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Request was '%s' (%d) and not OK (200). Body:\n%s\nurl: %s", resp.Status, resp.StatusCode, string(d), uri)
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func unzipFileAsData(f *zip.File) ([]byte, error) {
	r, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func zipExtract(d []byte) ([]*ReplitFile, error) {
	var res []*ReplitFile
	f := bytes.NewReader(d)
	fsize := int64(len(d))
	zr, err := zip.NewReader(f, fsize)
	if err != nil {
		return nil, err
	}
	for _, fi := range zr.File {
		if fi.FileInfo().IsDir() {
			continue
		}
		d, err := unzipFileAsData(fi)
		if err != nil {
			return nil, err
		}
		rf := &ReplitFile{
			name: fi.FileInfo().Name(),
			data: string(d),
		}
		res = append(res, rf)
	}
	return res, nil
}

func downloadAndCacheReplit(c *ReplitCache, uri string) (*Replit, bool, error) {
	fullURL := uri + ".zip"
	d, err := httpGet(fullURL)
	if err != nil {
		return nil, false, err
	}
	files, err := zipExtract(d)
	if err != nil {
		return nil, false, err
	}
	r := &Replit{
		url:   uri,
		files: files,
	}
	isNew, err := c.Add(r)
	return r, isNew, err
}
