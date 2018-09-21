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

	sha1HexCached string
}

// Sha1Hex returns sha1 of the content
func (f *FileContent) Sha1Hex() string {
	if f.sha1HexCached == "" {
		f.sha1HexCached = u.Sha1HexOfBytes(f.Content)
	}
	return f.sha1HexCached
}

var (
	filePathToFileContent = make(map[string]*FileContent)
	muFileCache           sync.Mutex
)

func cacheFile(path string, info os.FileInfo) (*FileContent, error) {
	d, err := common.ReadFileNormalized(path)
	if err != nil {
		return nil, err
	}
	s := string(d)
	lines := strings.Split(s, "\n")
	fc := &FileContent{
		Path:    path,
		Size:    info.Size(),
		ModTime: info.ModTime(),
		Content: d,
		Lines:   lines,
	}

	muFileCache.Lock()
	filePathToFileContent[path] = fc
	muFileCache.Unlock()

	return fc, nil
}

func cacheFileIfChanged(path string, info os.FileInfo) (*FileContent, error) {
	var err error
	if info == nil {
		info, err = os.Stat(path)
		if err != nil {
			return nil, err
		}
	}
	muFileCache.Lock()
	fc, ok := filePathToFileContent[path]
	muFileCache.Unlock()

	if !ok || fc.Size != info.Size() || fc.ModTime != info.ModTime() {
		return cacheFile(path, info)
	}
	return fc, nil
}

func loadFileCached(path string) (*FileContent, error) {
	return cacheFileIfChanged(path, nil)
}

func cacheFilesInDir(dir string) error {
	timeStart := time.Now()
	defer func() {
		fmt.Printf("cacheFilesInDir '%s' took %s\n", dir, time.Since(timeStart))
	}()
	res := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		_, err = cacheFileIfChanged(path, info)
		return err
	})
	return res
}
