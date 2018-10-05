// no playground
package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

// :show start
func main() {
	dir := "."
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("ioutil.ReadDir('%s') failed with '%s'\n", dir, err)
	}
	for i, fi := range fileInfos {
		if i < 4 {
			fmt.Printf("Path: %s, is dir: %v, size: %d bytes\n", fi.Name(), fi.IsDir(), fi.Size())
		}
	}
}

// :show end
