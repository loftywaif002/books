// no output
package main

import (
	"io"
	"log"
	"os"
)

// :show start
// CopyFile copies a src file to dst
func CopyFile(dst, src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	_, err = io.Copy(dstFile, srcFile)
	err2 := dstFile.Close()
	if err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		// delete the destination if copy failed
		os.Remove(dst)
	}
	return err
}

// :show end

func main() {
	src := "file_copy.go"
	dst := "file_copy_copy.go"
	err := CopyFile(dst, src)
	if err != nil {
		log.Fatalf("CopyFile failed with '%s'\n", err)
	}
	os.Remove(dst)
}
