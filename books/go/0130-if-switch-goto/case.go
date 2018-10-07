package main

import "fmt"

func main() {
	// :show start
	a := 1
	switch a {
	case 1, 3:
		fmt.Printf("a is 1 or 3\n")
	case 2:
		fmt.Printf("a is 2\n")
	default:
		fmt.Printf("default: a is %d\n", a)
	}
	// :show end
}
