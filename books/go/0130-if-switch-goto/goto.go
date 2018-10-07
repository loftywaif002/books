package main

import "fmt"

// :show start
func printIsOdd(n int) {
	if n%2 == 1 {
		goto isOdd
	}
	fmt.Printf("%d is even\n", n)
	return

isOdd:
	fmt.Printf("%d is odd\n", n)
}

// :show end

func main() {
	// :show start
	printIsOdd(5)
	printIsOdd(18)
	// :show end
}
