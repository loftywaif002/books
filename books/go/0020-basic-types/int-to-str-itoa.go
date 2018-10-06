package main

import (
	"fmt"
	"strconv"
)

func main() {
	// :show start
	var i1 int = -38
	fmt.Printf("i1: %s\n", strconv.Itoa(i1))

	var i2 int32 = 148
	fmt.Printf("i2: %s\n", strconv.Itoa(int(i2)))
	// :show end
}
