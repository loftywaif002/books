package main

import (
	"fmt"
)

// :show start
func check(n int) {
	switch {
	case n > 0 && n%3 == 0:
		fmt.Printf("n is %d, divisible by 3\n", n)
	case n >= 4:
		fmt.Printf("n is %d (>= 4)\n", n)
	default:
		fmt.Printf("default: n is %d\n", n)
	}
}

// :show end
func main() {
	// :show start
	check(3)
	check(4)
	check(6)
	check(1)
	// :show end
}
