// allow error
package main

import (
	"fmt"
	"reflect"
)

// :show start

func getIntValue(v interface{}) {
	var reflectValue = reflect.ValueOf(v)
	n := reflectValue.Int()
	fmt.Printf("Int value is: %d\n", n)
}

// :show end

func main() {
	// :show start
	getIntValue(3)
	getIntValue(int8(4))
	getIntValue("")
	// :show end
}
