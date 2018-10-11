package main

import (
	"fmt"
	"log"
)

func main() {
	// :show start
	// extract int and float from a string
	s := "48 123.45"
	var f float64
	var i int
	nParsed, err := fmt.Sscanf(s, "%d %f", &i, &f)
	if err != nil {
		log.Fatalf("first fmt.Sscanf failed with %s\n", err)
	}
	fmt.Printf("i: %d, f: %f, extracted %d values\n", i, f, nParsed)

	var i2 int
	_, err = fmt.Sscanf(s, "%d %f %d", &i, &f, &i2)
	if err != nil {
		fmt.Printf("second fmt.Sscanf failed with %s\n", err)
	}

	// :show end
}
