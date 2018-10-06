// no output
package main

type Person struct {
	FirstName string
	LastName  string
}

func main() {
	// :show start
	// Basic variable declaration. Declares a variable of type specified on the right.
	// The variable is initialized to the zero value of the respective type.
	var x int
	var s string
	var p Person // Assuming type Person struct {}

	// Assignment of a value to a variable
	x = 3

	// Short declaration using := infers the type
	y := 4

	u := int64(100)    // declare variable of type int64 and init with 100
	var u2 int64 = 100 // declare variable of type int64 and init with 100
	// :show end

	// silence compiler error about unused variables
	_, _, _, _, _, _ = x, s, p, y, u, u2
}
