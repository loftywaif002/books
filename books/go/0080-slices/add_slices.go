package main

import "fmt"

func main() {
	// :show start
	a := []string{"!"}
	a2 := []string{"Hello", "world"}
	a = append(a, a2...)
	fmt.Printf("a: %#v\n", a)
	// :show end
}
