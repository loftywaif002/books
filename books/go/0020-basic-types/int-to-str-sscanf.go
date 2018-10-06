package main

import (
	"fmt"
	"log"
)

func main() {
	// :show start
	s := "348"
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil {
		log.Fatalf("fmt.Sscanf failed with '%s'\n", err)
	}
	fmt.Printf("i1: %d\n", i)
	// :show end
}
