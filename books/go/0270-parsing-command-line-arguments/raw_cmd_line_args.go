// :run go run $file -echo echo-arg additional arg
// no playground
package main

import (
	"fmt"
	"os"
)

// :show start
func main() {
	fmt.Printf("Name of executable: '%s'\n", os.Args[0])
	args := os.Args[1:]
	for i, arg := range args {
		fmt.Printf("Arg %d, value: '%s'\n", i, arg)
	}
}

// :show end
