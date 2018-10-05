// allow error
package main

import "fmt"

// :show start
func logNum(i int) {
	fmt.Printf("Num %d\n", i)
}

func main() {
	defer logNum(1)
	fmt.Println("First main statement")
	defer logNum(2)
	defer logNum(3)
	panic("panic occurred")

	fmt.Println("Last main statement") // not printed

	// not deferred since execution flow never reaches this line
	defer logNum(3)
}

// :show end
