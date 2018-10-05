// allow error
package main

import "fmt"

// :show start
func foo() {
	defer fmt.Println("Exiting foo")
	panic("bar")
}

func main() {
	defer fmt.Println("Exiting main")
	foo()
}

// :show end
