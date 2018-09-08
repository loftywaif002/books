package main

import (
	"fmt"
	"log"
)

func main() {
	// :show start
	s := "348"
	var i1 int
	_, err := fmt.Sscanf(s, "%d", &f)
	if err != nil {
		log.Fatalf("fmt.Sscanf failed with '%s'\n", err)
	}
	fmt.Printf("i1: %d\n", i1)
	// :show end
}
