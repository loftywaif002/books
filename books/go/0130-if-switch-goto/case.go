package main

import "fmt"

func main() {
	// :show start
	a := 1
	switch a {
	case 1:
		fmt.Printf("a is 1\n")
	case 2:
		fmt.Printf("a is 2\n")
	default:
		fmt.Printf("a is not 1 or 2 but %d\n", a)
	}
	// :show end
}
