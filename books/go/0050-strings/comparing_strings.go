package main

import (
	"fmt"
)

func main() {
	// :show start
	s1 := "string one"
	s2 := "string two"

	if s1 == s2 {
		fmt.Printf("s1 is equal to s2\n")
	} else {
		fmt.Printf("s1 is not equal to s2\n")
	}

	if s1 == s1 {
		fmt.Printf("s1 is equal to s1\n")
	} else {
		fmt.Printf("inconcivable! s1 is not equal to itself\n")
	}

	if s1 > s2 {
		fmt.Printf("s1 is > than s2\n")
	} else {
		fmt.Printf("s1 is not > than s2\n")
	}

	if s1 < s2 {
		fmt.Printf("s1 is < than s2\n")
	} else {
		fmt.Printf("s1 is not < than s2\n")
	}
	// :show end
}
