package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/essentialbooks/books/pkg/common"
	"github.com/kjk/notionapi"
	"github.com/kjk/u"
)

// SourceFile represents source file present in the repository
// and embedded via https://www.onlinetool.io/gitoembed/
type SourceFile struct {
	EmbedURL string

	// name of the file
	FileName string
	// full path of the file
	Path string

	FileExists bool

	// content of the file after filtering
	Lines         []string
	cachedData    []byte
	cachedSha1Hex string

	// output of running a file
	Output []byte
}

// Data returns content of the file
func (f *SourceFile) Data() []byte {
	if len(f.cachedData) == 0 {
		s := strings.Join(f.Lines, "\n")
		f.cachedData = []byte(s)
	}
	return f.cachedData
}

// RealSha1Hex returns hex version of sha1 of file content
func (f *SourceFile) RealSha1Hex() string {
	if f.cachedSha1Hex == "" {
		f.cachedSha1Hex = u.Sha1HexOfBytes(f.Data())
	}
	return f.cachedSha1Hex
}

// https://www.onlinetool.io/gitoembed/widget?url=https%3A%2F%2Fgithub.com%2Fessentialbooks%2Fbooks%2Fblob%2Fmaster%2Fbooks%2Fgo%2F0020-basic-types%2Fbooleans.go
// to:
// books/go/0020-basic-types/booleans.go
// returns empty string if doesn't conform to what we expect
func gitoembedToRelativePath(uri string) string {
	parsed, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	switch parsed.Host {
	case "www.onlinetool.io", "onlinetool.io":
		// do nothing
	default:
		return ""
	}
	path := parsed.Path
	if path != "/gitoembed/widget" {
		return ""
	}
	uri = parsed.Query().Get("url")
	// https://github.com/essentialbooks/books/blob/master/books/go/0020-basic-types/booleans.go
	parsed, err = url.Parse(uri)
	if parsed.Host != "github.com" {
		return ""
	}
	path = strings.TrimPrefix(parsed.Path, "/essentialbooks/books/")
	if path == parsed.Path {
		return ""
	}
	// blob/master/books/go/0020-basic-types/booleans.go
	path = strings.TrimPrefix(path, "blob/")
	// master/books/go/0020-basic-types/booleans.go
	// those are branch names. Should I just strip first 2 elements from the path?
	path = strings.TrimPrefix(path, "master/")
	path = strings.TrimPrefix(path, "notion/")
	// books/go/0020-basic-types/booleans.go
	return path
}

// we don't want to show our // :show annotations in snippets
func removeAnnotationLines(lines []string) []string {
	var res []string
	prevWasEmpty := false
	for _, l := range lines {
		if strings.Contains(l, "// :show ") {
			continue
		}
		if len(l) == 0 && prevWasEmpty {
			continue
		}
		prevWasEmpty = len(l) == 0
		res = append(res, l)
	}
	return res
}

func readFilteredSourceFile(path string) ([]string, error) {
	d, err := common.ReadFileNormalized(path)
	if err != nil {
		return nil, err
	}
	lines := dataToLines(d)
	lines = removeAnnotationLines(lines)
	return lines, nil
}

func extractSourceFiles(p *Page) {
	wd, err := os.Getwd()
	panicIfErr(err)
	page := p.NotionPage
	for _, block := range page.Root.Content {
		if block.Type != notionapi.BlockEmbed {
			continue
		}
		uri := block.FormatEmbed.DisplaySource
		f := &SourceFile{
			EmbedURL: uri,
		}
		p.SourceFiles = append(p.SourceFiles, f)
		relativePath := gitoembedToRelativePath(uri)
		if relativePath == "" {
			fmt.Printf("Couldn't parse embed uri '%s'\n", uri)
			continue
		}
		// fmt.Printf("Embed uri: %s, relativePath: %s\n", uri, relativePath)
		f.FileName = filepath.Base(relativePath)
		f.Path = filepath.Join(wd, relativePath)
		f.Lines, err = readFilteredSourceFile(f.Path)
		if err != nil {
			fmt.Printf("Failed to read '%s' extracted from '%s', error: %s\n", f.Path, uri, err)
			continue
		}
		f.FileExists = true
	}
}
