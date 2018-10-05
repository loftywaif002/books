// allow error
package main

import (
	"fmt"
)

// :show start
func panicOnInvalidConversion(iv interface{}) {
	v := iv.(int)
	fmt.Printf("v is int of value: %d\n", v)
}

func main() {
	panicOnInvalidConversion("string")
}

// :show end
