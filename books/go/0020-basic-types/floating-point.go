package main

import (
	"fmt"
	"strconv"
)

func main() {
	// :show start
	var f32 float32 = 1.3
	bitSize := 32
	s1 := strconv.FormatFloat(float64(f32), 'E', -1, bitSize)
	fmt.Printf("f32: %s\n", s1)

	var f64 float64 = 8.1234
	bitSize = 64
	s2 := strconv.FormatFloat(f64, 'e', -1, bitSize)
	fmt.Printf("f64: %s\n", s2)
	// :show end
}
