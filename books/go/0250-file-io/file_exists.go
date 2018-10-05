// no playground
package main

import (
	"fmt"
	"os"
)

// :show start

// IsPathxists returns true if a given path exists, false if it doesn't.
// It might return an error if e.g. file exists but you don't have
// access
func IsPathxists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	// error other than not existing e.g. permission denied
	return false, err
}

func printExists(path string) {
	exists, err := IsPathxists(path)
	if err == nil {
		fmt.Printf("File '%s' exists: %v\n", path, exists)
	} else {
		fmt.Printf("IsFileExists('%s') failed with '%s'\n", path, err)
	}
}
func main() {
	printExists("file_exists.go")
	printExists("non-existent-file.txt")
}

// :show end
