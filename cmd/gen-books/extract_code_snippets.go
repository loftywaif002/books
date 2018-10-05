package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// TODO: rename to: code_extract_snippets.go

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

func extractCodeSnippets(lines []string) ([]string, error) {
	var res [][]string
	var curr []string
	inShow := false
	for _, line := range lines {
		if isShowStart(line) {
			if inShow {
				return nil, fmt.Errorf("consequitive '%s' lines", showStartLine)
			}
			inShow = true
			continue
		}
		if isShowEnd(line) {
			if !inShow {
				return nil, fmt.Errorf("'%s' without start line", showEndLine)
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

// finds ":run ${cmd}" directive embedded in the file
// and returns ${cmd} part or empty string if not found
func extractRunCmd(lines []string) (string, []string) {
	for i, line := range lines {
		if idx := strings.Index(line, ":run "); idx != -1 {
			s := line[idx+len(":run "):]
			s = strings.TrimSpace(s)
			lines = append(lines[:i], lines[i+1:]...)
			return s, lines
		}
	}
	return "", lines
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
