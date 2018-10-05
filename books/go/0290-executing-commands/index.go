// no playground
package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	// :show start
	cmd := exec.Command("go", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.CombinedOutput() failed with '%s'\n", err)
	}
	fmt.Printf("Output:\n%s\n", string(out))
	// :show end
}
