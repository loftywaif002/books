package main

import "fmt"

func main() {
	// :show start
	m := make(map[string]int)
	m["number3"] = 3

	k := "number3"
	if n, ok := m[k]; ok {
		fmt.Printf("value for key '%s' is %d\n", k, n)
	} else {
		fmt.Printf("key '%s' doesn't exist in map m\n", k)
	}

	k = "number4"
	if n, ok := m[k]; ok {
		fmt.Printf("value for key '%s' is %d\n", k, n)
	} else {
		fmt.Printf("key '%s' doesn't exist in map m\n", k)
	}
	// :show end
}
