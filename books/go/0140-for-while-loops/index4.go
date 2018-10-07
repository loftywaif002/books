package main

import "fmt"

func main() {
	// :show start
	for i := 0; i < 5; {
		fmt.Printf("i: %d\n", i)
		i += 2
	}
	// :show end
}
