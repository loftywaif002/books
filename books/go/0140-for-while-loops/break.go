package main

import "fmt"

func main() {
	// show start
	i := 0
	for {
		i++
		if i > 2 {
			break
		}
		fmt.Printf("i: %d\n", i)
	}
	// show end
}
