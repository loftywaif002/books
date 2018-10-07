package main

import "fmt"

func main() {
	// :show start
	i := 0
	for ; ; i += 2 {
		fmt.Printf("i: %d\n", i)
		if i >= 5 {
			break
		}
	}
	// :show end
}
