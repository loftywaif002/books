package main

import (
	"fmt"
	"strings"
)

func main() {
	// :show start
	s := "this,. is,. a,. string"
	a := strings.Split(s, ",.")
	fmt.Printf("a: %#v\n", a)
	// :show end
}
