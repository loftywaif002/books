package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

/*
We scan all .go files under books/go, extract external dependencies
and run go get -u on this dependency, to make sure we test against the latest.

This is meant to be run once on startup.
*/

var (
	updatedImports = map[string]bool{}
)

func shouldUpdatePackage(name string) bool {
	// we don't want to update built-in packages (like "fmt")
	// we assume that non-built-in packages must include domain name
	// in the name and domain name must contain ".", while built-in package
	// will not
	return strings.Contains(name, ".")
}

func updateImport(importSpec *ast.ImportSpec) {
	if importSpec.Path == nil {
		return
	}
	path := importSpec.Path
	name := path.Value
	// at this point name is sth. like: "fmt"
	// strip '"' from both ends
	name = name[1 : len(name)-1]
	// we've already updated this import
	if _, ok := updatedImports[name]; ok {
		return
	}
	updatedImports[name] = true
	if !shouldUpdatePackage(name) {
		//fmt.Printf("not doing go get -u %s\n", name)
		return
	}
	fmt.Printf("running go get -u %s\n", name)
	cmd := exec.Command("go", "get", "-u", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("go get failed with error: '%s', Output: '%s'\n", err, string(out))
	}
}

func updateGoDepsInFile(path string) {
	//fmt.Printf("File: %s\n", path)
	// Inspect the AST and print all identifiers and literals.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		panic(err)
	}

	astInspect := func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			updateImport(x)
		}
		return true
	}
	ast.Inspect(f, astInspect)
}

func updateGoDeps() {
	timeStart := time.Now()
	fmt.Printf("udpateGoDeps() start\n")
	dir := filepath.Join("books", "go")
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".go" {
			return nil
		}
		updateGoDepsInFile(path)
		return nil
	}
	filepath.Walk(dir, walkFunc)
	fmt.Printf("updateGoDeps() finished in %s\n", time.Since(timeStart))
}
