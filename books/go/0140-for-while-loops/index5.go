package main

import "fmt"

func main() {
	// :show start
	i := 0
	for {
		fmt.Printf("i: %d\n", i)
		i += 2
		if i >= 5 {
			break
		}
	}

	// :show end
}
