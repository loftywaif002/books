package main

import (
	"fmt"
	"log"
	"strconv"
)

func main() {
	// :show start
	s := "1.2341"
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatalf("strconv.ParseFloat() failed with '%s'\n", err)
	}
	fmt.Printf("f64: %f\n", f64)
	// :show end
}
