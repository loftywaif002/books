// no playground
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func exitIfError(err error) {
	if err != nil {
		fmt.Printf("Error: '%s'\n", err)
		os.Exit(1)
	}
}

func main() {
	// :show start
	path := "index.go"
	f, err := os.Open(path)
	exitIfError(err)
	defer f.Close()

	d, err := ioutil.ReadAll(f)
	exitIfError(err)

	lines := bytes.Split(d, []byte{'\n'})
	fmt.Printf("File %s has %d lines\n", path, len(lines))
	// :show end
}
