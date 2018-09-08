package main

import (
	"fmt"
	"log"
	"strconv"
)

func main() {
	// :show start
	s := "-48"
	i1, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("strconv.Atoi() failed with %s\n", err)
	}
	fmt.Printf("i1: %d\n", i1)
	// :show end
}
