// no playground
package main

import (
	"log"
	"os/exec"
)

func main() {
	// :show start
	cmd := exec.Command("go", "version")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("first cmd.CombintedOutput() failed with '%s'\n", err)
	}

	_, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("second cmd.CombintedOutput() failed with '%s'\n", err)
	}
	// :show end
}
