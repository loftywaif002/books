// no playground
package main

import (
	"fmt"
	"log"
	"os"
)

// :show start

func main() {
	st, err := os.Stat("file_info.go")
	if err != nil {
		log.Fatalf("GetFileSize failed with '%s'\n", err)
	}
	fmt.Printf(`Name: %s
Size: %d
IsDir: %v
Mode: %x
ModTime: %s
OS info: %T
`, st.Name(), st.Size(), st.IsDir(), st.Mode, st.ModTime(), st.Sys())
}

// :show end
