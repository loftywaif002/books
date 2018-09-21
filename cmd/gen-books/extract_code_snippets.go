package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/shlex"
	"github.com/kjk/u"
)

const (
	showStartLine = "// :show start"
	showEndLine   = "// :show end"
	// if false, we separate code snippet and output
	// with **Output** paragraph
	compactOutput = true
)

func isShowStart(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == showStartLine
}

func isShowEnd(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == showEndLine
}

func extractCodeSnippets(path string) ([]string, error) {
	//fmt.Printf("extractCodeSnippets: %s\n", path)
	fc, err := loadFileCached(path)
	if err != nil {
		return nil, err
	}
	lines := fc.Lines
	var res [][]string
	var curr []string
	inShow := false
	for _, line := range lines {
		if isShowStart(line) {
			if inShow {
				return nil, fmt.Errorf("file '%s': consequitive '%s' lines", path, showStartLine)
			}
			inShow = true
			continue
		}
		if isShowEnd(line) {
			if !inShow {
				return nil, fmt.Errorf("file '%s': '%s' without start line", path, showEndLine)
			}
			inShow = false
			if len(curr) > 0 {
				res = append(res, curr)
			}
			curr = nil
			continue
		}
		if inShow {
			curr = append(curr, line)
		}
	}
	// if there are no show: markings, assume we want to show the whole file
	if len(res) == 0 {
		return trimEmptyLines(lines), nil
	}
	var all []string
	for _, lines := range res {
		shiftLines(lines)
		all = append(all, lines...)
		// add a separation line between show sections.
		// should be the right thing more often than not
		all = append(all, "")
	}
	return trimEmptyLines(all), nil
}

func getLangFromFileExt(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".go":
		return "go"
	case ".json":
		return "js"
	case ".csv":
		// note: chroma doesn't have csv lexer
		return "text"
	case ".yml":
		return "yaml"
	}
	fmt.Printf("Couldn't deduce language from file name '%s'\n", fileName)
	// TODO: more languages
	return ""
}

// convert local path like books/go/foo.go into path to the file in a github repo
func getGitHubPathForFile(path string) string {
	return "https://github.com/essentialbooks/books/blob/master/" + toUnixPath(path)
}

// ${baseDir} is books/go/
// loads a source file whose name is in ${line} and
func extractCodeSnippetsAsMarkdownLines(baseDir string, line string) ([]string, error) {
	// line is:
	// @file ${fileName} [output]
	directive, err := parseFileDirective(line)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(baseDir, directive.FileName)
	if !fileExists(path) {
		return nil, fmt.Errorf("no file '%s' in line '%s'", path, line)
	}
	lines, err := extractCodeSnippets(path)
	if err != nil {
		return nil, err
	}
	lang := getLangFromFileExt(path)
	sep := "|"
	u.PanicIf(strings.Contains(lang, sep), "lang ('%s') contains '%s'", lang, sep)
	u.PanicIf(strings.Contains(path, sep), "path ('%s') contains '%s'", path, sep)
	// this line is parsed in parseCodeBlockInfo
	s := fmt.Sprintf("%s|github|%s", lang, getGitHubPathForFile(path))
	if directive.GoPlaygroundID != "" {
		// alternative would be https://play.golang.org/p/ + ${id}
		uri := "https://goplay.space/#" + directive.GoPlaygroundID
		s += "|playground|" + uri
	}
	if directive.LineLimit != 0 {
		n := directive.LineLimit
		if n < len(lines) {
			lines = lines[:n]
		}
	}
	res := []string{"```" + s}
	res = append(res, lines...)
	res = append(res, "```")

	if !directive.WithOutput {
		return res, nil
	}

	out, err := getCachedOutput(path, directive.AllowError)
	if err != nil {
		fmt.Printf("getCachedOutput('%s'): error '%s', output: '%s'\n", path, err, out)
		maybePanicIfErr(err)
		return res, err
	}

	if compactOutput {
		res = append(res, "")
		res = append(res, "```output")
	} else {
		res = append(res, "")
		res = append(res, "**Output**:")
		res = append(res, "")
		res = append(res, "```text")
	}
	lines = strings.Split(out, "\n")
	lines = trimEmptyLines(lines)
	res = append(res, lines...)
	res = append(res, "```")
	return res, nil
}

// runs `go run ${path}` and returns captured output`
func getGoOutput(path string) (string, error) {
	dir, fileName := filepath.Split(path)
	cmd := exec.Command("go", "run", fileName)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func getRunCmdOutput(path string, runCmd string) (string, error) {
	parts, err := shlex.Split(runCmd)
	maybePanicIfErr(err)
	if err != nil {
		return "", err
	}
	exeName := parts[0]
	parts = parts[1:]
	var parts2 []string
	srcDir, srcFileName := filepath.Split(path)

	// remove empty lines and replace variables
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		switch part {
		case "$file":
			part = srcFileName
		}
		parts2 = append(parts2, part)
	}
	//fmt.Printf("getRunCmdOutput: running '%s' with args '%#v'\n", exeName, parts2)
	cmd := exec.Command(exeName, parts2...)
	cmd.Dir = srcDir
	out, err := cmd.CombinedOutput()
	//fmt.Printf("getRunCmdOutput: out:\n%s\n", string(out))
	return string(out), err
}

// finds ":run ${cmd}" directive embedded in the file
// and returns ${cmd} part or empty string if not found
func findRunCmd(lines []string) string {
	for _, line := range lines {
		if idx := strings.Index(line, ":run "); idx != -1 {
			s := line[idx+len(":run "):]
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func stripCurrentPathFromOutput(s string) string {
	path, err := filepath.Abs(".")
	u.PanicIfErr(err)
	return strings.Replace(s, path, "", -1)
}

// it executes a code file and captures the output
// optional runCmd says
func getOutput(path string) (string, error) {
	fc, err := loadFileCached(path)
	if err != nil {
		return "", err
	}
	if runCmd := findRunCmd(fc.Lines); runCmd != "" {
		//fmt.Printf("Found :run cmd '%s' in '%s'\n", runCmd, path)
		s, err := getRunCmdOutput(path, runCmd)
		return stripCurrentPathFromOutput(s), err
	}

	// do default
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".go" {
		s, err := getGoOutput(path)
		return stripCurrentPathFromOutput(s), err
	}
	return "", fmt.Errorf("getOutpu(%s): files with extension '%s' are not supported", path, ext)
}
