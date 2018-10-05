// :run go run $file -echo echo-arg additional arg
// no playground
package main

import (
	"flag"
	"fmt"
	"os"
)

// :show start
var (
	flgHelp bool
	flgEcho string
)

func parseCmdLineFlags() {
	flag.BoolVar(&flgHelp, "help", false, "if true, show help")
	flag.StringVar(&flgEcho, "echo", "", "")
	flag.Parse()
}

func main() {
	parseCmdLineFlags()
	if flgHelp {
		flag.Usage()
		os.Exit(0)
	}
	if flgEcho != "" {
		fmt.Printf("flag -echo: '%s'\n", flgEcho)
	}

	remainingArgs := flag.Args()
	for _, arg := range remainingArgs {
		fmt.Printf("Remainig arg: '%s'\n", arg)
	}
}

// :show end
