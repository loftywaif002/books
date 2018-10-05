// no playground
package main

import (
	"fmt"
	"log"
	"os"
)

// :show start

// GetFileSize returns file size or error if e.g. file doesn't exist
func GetFileSize(path string) (int64, error) {
	st, err := os.Lstat(path)
	if err != nil {
		return -1, err
	}
	return st.Size(), nil
}

func main() {
	path := "file_size.go"
	size, err := GetFileSize(path)
	if err != nil {
		log.Fatalf("GetFileSize failed with '%s'\n", err)
	}
	fmt.Printf("File %s is %d bytes in size\n", path, size)
}

// :show end
