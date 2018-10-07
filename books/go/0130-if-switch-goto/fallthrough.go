package main

import "fmt"

func main() {
	// :show start
	a := 1
	switch a {
	case 1:
		fmt.Printf("case 1\n")
		fallthrough
	case 2:
		fmt.Printf("caes 2\n")
	}
	// :show end
}
