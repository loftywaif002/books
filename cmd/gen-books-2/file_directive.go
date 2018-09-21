package main

import (
	"fmt"
	"strconv"
	"strings"
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
	parts := strings.Split(s, ",")
	for _, s := range parts {
		s = strings.TrimSpace(s)
		switch s {
		case "no output":
			res.NoOutput = true
		case "no playground":
			res.NoPlayground = true
		case "allow error":
			res.AllowError = true
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
		}
	}
	// TODO: set NoPlayground for non-Go files
	return res, nil
}
