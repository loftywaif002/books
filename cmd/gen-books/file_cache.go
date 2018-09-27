package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/essentialbooks/books/pkg/common"
	"github.com/kjk/u"
)

// FileContent describes a file and its content
type FileContent struct {
	Path    string
	Size    int64
	ModTime time.Time
	Content []byte
	Lines   []string

	directive     *FileDirective
	sha1HexCached string
}

// Sha1Hex returns sha1 of the content
func (f *FileContent) Sha1Hex() string {
	if f.sha1HexCached == "" {
		f.sha1HexCached = u.Sha1HexOfBytes(f.Content)
	}
	return f.sha1HexCached
}

// FileCache is a cache of files based on their paths
type FileCache struct {
	pathToFileContent map[string]*FileContent
	mu                sync.Mutex
}

// NewFileCache creates FileCache
func NewFileCache() *FileCache {
	return &FileCache{
		pathToFileContent: map[string]*FileContent{},
	}
}

func (c *FileCache) cacheFile(path string, info os.FileInfo) (*FileContent, error) {
	d, err := common.ReadFileNormalized(path)
	if err != nil {
		return nil, err
	}
	s := string(d)
	lines := strings.Split(s, "\n")
	var directive *FileDirective
	if len(lines) > 0 {
		directive, err = parseFileDirective(lines[0])
		panicIfErr(err)
		if directive != nil {
			lines = lines[1:]
		}
	}
	fc := &FileContent{
		Path:      path,
		Size:      info.Size(),
		ModTime:   info.ModTime(),
		Content:   d,
		Lines:     lines,
		directive: directive,
	}

	c.mu.Lock()
	c.pathToFileContent[path] = fc
	c.mu.Unlock()

	return fc, nil
}

func (c *FileCache) cacheFileIfChanged(path string, info os.FileInfo) (*FileContent, error) {
	var err error
	if info == nil {
		info, err = os.Stat(path)
		if err != nil {
			return nil, err
		}
	}
	c.mu.Lock()
	fc, ok := c.pathToFileContent[path]
	c.mu.Unlock()

	if !ok || fc.Size != info.Size() || fc.ModTime != info.ModTime() {
		return c.cacheFile(path, info)
	}
	return fc, nil
}

func (c *FileCache) loadFileCached(path string) (*FileContent, error) {
	return c.cacheFileIfChanged(path, nil)
}

func cacheFilesInDir(dir string, allowFile func(string) bool) (*FileCache, error) {
	timeStart := time.Now()
	cache := NewFileCache()
	defer func() {
		fmt.Printf("cacheFilesInDir '%s' took %s\n", dir, time.Since(timeStart))
	}()
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		allow := true
		if allowFile != nil {
			allow = allowFile(path)
		}
		if allow {
			_, err = cache.cacheFileIfChanged(path, info)
		}
		return err
	})
	return cache, err
}
