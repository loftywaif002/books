// allow error
package main

import (
	"fmt"
)

// :show start
func printTypeAndValue(iv interface{}) {
	if v, ok := iv.(string); ok {
		fmt.Printf("iv is of type string and has value '%s'\n", v)
		return
	}
	if v, ok := iv.(int); ok {
		fmt.Printf("iv is of type int and has value '%d'\n", v)
		return
	}
	if v, ok := iv.(*int); ok {
		fmt.Printf("iv is of type *int and has value '%s'\n", v)
		return
	}
}

func panicOnInvalidConversion() {
	var iv interface{} = "string"

	v := iv.(int)
	fmt.Printf("v is int of value: %d\n", v)
}

func main() {
	// pass a string
	printTypeAndValue("string")
	i := 5
	// pass an int
	printTypeAndValue(i)
	// pass a pointer to int i.e. *int
	printTypeAndValue(&i)

	panicOnInvalidConversion()
}

// :show end
