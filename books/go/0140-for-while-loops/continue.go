package main

import "fmt"

func main() {
	// :show start
	for i := 0; i < 4; i++ {
		if i < 2 {
			continue
		}
		fmt.Printf("i: %d\n", i)
	}
	// :show end
}
