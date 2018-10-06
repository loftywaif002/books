// no output
package main

func main() {
	// :show start
	// You can declare multiple variables of the same type in one line
	var a, b, c string

	var d, e string = "Hello", "world!"

	// You can also use short declaration to assign multiple variables
	x, y, z := 1, 2, 3

	foo, bar := 4, "stack" // `foo` is type `int`, `bar` is type `string`
	// :show end

	// silence compiler error about unused variables
	_, _, _, _, _, _, _, _, _, _ = a, b, c, d, e, x, y, z, foo, bar
}
