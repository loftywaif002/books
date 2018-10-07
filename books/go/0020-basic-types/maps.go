package main

import "fmt"

func main() {
	// :show start
	m := make(map[string]int)
	m["number3"] = 3

	checkKey := func(k string) {
		if n, ok := m[k]; ok {
			fmt.Printf("value for key '%s' is %d\n", k, n)
		} else {
			fmt.Printf("key '%s' doesn't exist in map\n", k)
		}
	}

	// get value for a key that exists
	checkKey("number3")

	// get value for a key that doesn't exist
	checkKey("number4")

	// remove a key
	delete(m, "number3")
	fmt.Printf("deleted key 'number3\n")
	checkKey("number3") // and now it doesn't exist
	// :show end
}
