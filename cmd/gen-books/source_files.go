package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/essentialbooks/books/pkg/common"
	"github.com/kjk/notionapi"
)

/*
FileDirective describes reulst of parsing a line like:
// no output, no playground
*/
type FileDirective struct {
	NoOutput     bool // "no output"
	AllowError   bool // "allow error"
	LineLimit    int  // limit ${n}
	NoPlayground bool // no playground
}

/* Parses a line like:
// no output, no playground, line ${n}, allow error
*/
func parseFileDirective(line string) (*FileDirective, error) {
	line = strings.TrimSpace(line)
	s := strings.TrimSuffix(line, "//")
	// doesn't start with a comment, so is not a file directive
	if s == line {
		return nil, nil
	}
	res := &FileDirective{}
	hasInfo := false
	parts := strings.Split(s, ",")
	for _, s := range parts {
		s = strings.TrimSpace(s)
		switch s {
		case "no output":
			res.NoOutput = true
			hasInfo = true
		case "no playground", "noplayground":
			res.NoPlayground = true
			hasInfo = true
		case "allow error":
			res.AllowError = true
			hasInfo = true
		default:
			rest := strings.TrimPrefix(s, "line ")
			if rest == s {
				return nil, fmt.Errorf("parseFileDirective: invalid line '%s'", line)
			}
			n, err := strconv.Atoi(rest)
			if err != nil {
				return nil, fmt.Errorf("parseFileDirective: invalid line '%s'", line)
			}
			res.LineLimit = n
			hasInfo = true
		}
	}
	if !hasInfo {
		return nil, nil
	}
	return res, nil
}

func extractFileDirective(lines []string) (*FileDirective, []string, error) {
	directive, err := parseFileDirective(lines[0])
	if err != nil {
		return nil, nil, err
	}
	if directive == nil {
		return &FileDirective{}, lines, nil
	}
	return directive, lines[1:], nil
}

// SourceFile represents source file present in the repository
// and embedded via https://www.onlinetool.io/gitoembed/
type SourceFile struct {
	EmbedURL string

	// full path of the file
	Path string
	// name of the file
	FileName string
	// URL on GitHub for this file
	GitHubURL string
	// language of the file, detected from name
	Lang string

	// optional, ":run ${cmd}" extracted from file content
	RunCmd string

	// for Go files, this is playground id
	GoPlaygroundID string

	// optional, extracted from first line of the file
	// allows providing meta-data instruction for this file
	Directive *FileDirective

	// raw content of the file with line endings normalized to '\n'
	Data []byte

	LinesRaw []string // Data split into lines

	// LinesRaw after extracting directive, run cmd at the top
	// and removing :show annotation lines
	// This is the content sent to playgrounds
	LinesFiltered []string

	// the part that we want to show i.e. the parts inside
	// :show start, :show end blocks
	LinesCode []string

	// output of running a file
	Output string
}

// DataFiltered returns content of the file after filtering
func (f *SourceFile) DataFiltered() []byte {
	s := strings.Join(f.LinesFiltered, "\n")
	return []byte(s)
}

// DataCode returns part of the file tbat we want to show
func (f *SourceFile) DataCode() []byte {
	s := strings.Join(f.LinesCode, "\n")
	return []byte(s)
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

// convert local path like books/go/foo.go into path to the file in a github repo
func getGitHubPathForFile(path string) string {
	return "https://github.com/essentialbooks/books/blob/master/" + toUnixPath(path)
}

func setGoPlaygroundID(sf *SourceFile) error {
	if sf.Lang != "go" {
		return nil
	}
	if sf.Directive.NoPlayground {
		return nil
	}
	id, err := getSha1ToGoPlaygroundIDCached(sf.DataFiltered())
	if err != nil {
		return err
	}
	sf.GoPlaygroundID = id
	return nil
}

func loadSourceFile(path string) (*SourceFile, error) {
	data, err := common.ReadFileNormalized(path)
	if err != nil {
		return nil, err
	}
	name := filepath.Base(path)
	lang := getLangFromFileExt(filepath.Ext(path))
	gitHubURL := getGitHubPathForFile(path)
	sf := &SourceFile{
		Path:      path,
		FileName:  name,
		Data:      data,
		Lang:      lang,
		GitHubURL: gitHubURL,
	}
	sf.LinesRaw = dataToLines(sf.Data)
	lines := sf.LinesRaw
	sf.RunCmd, lines = extractRunCmd(lines)
	directive, lines, err := extractFileDirective(lines)
	if err != nil {
		fmt.Printf("loadSourceFile: extractFileDirective() of line '%s' failed with '%s'\n", sf.LinesRaw[0], err)
		panicIfErr(err)
	}
	sf.Directive = directive

	sf.LinesFiltered = removeAnnotationLines(lines)
	sf.LinesCode, err = extractCodeSnippets(lines)
	if err != nil {
		fmt.Printf("loadSourceFile('%s'): extractCodeSnippets() failed with '%s'\n", path, err)
		panicIfErr(err)
	}
	setGoPlaygroundID(sf)
	err = getOutputCached(sf)
	fmt.Printf("loadSourceFile('%s'), lang: '%s'\n", path, lang)
	return sf, nil
}

func extractSourceFiles(p *Page) {
	//wd, err := os.Getwd()
	//panicIfErr(err)
	page := p.NotionPage
	for _, block := range page.Root.Content {
		if block.Type != notionapi.BlockEmbed {
			continue
		}
		uri := block.FormatEmbed.DisplaySource
		relativePath := gitoembedToRelativePath(uri)
		if relativePath == "" {
			fmt.Printf("Couldn't parse embed uri '%s'\n", uri)
			continue
		}
		// fmt.Printf("Embed uri: %s, relativePath: %s\n", uri, relativePath)
		//path := filepath.Join(wd, relativePath)
		path := relativePath
		sf, err := loadSourceFile(path)
		if err != nil {
			fmt.Printf("extractSourceFiles: loadSourceFile('%s') (uri: '%s') failed with '%s'\n", path, uri, err)
			panicIfErr(err)
		}
		sf.EmbedURL = uri
		p.SourceFiles = append(p.SourceFiles, sf)
	}
}
