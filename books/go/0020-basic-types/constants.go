// no output
package main

// :show start
const (
	i  int = 32       // int constant
	s      = "string" // string constant
	i2     = 33       // untyped number constant
)

var (
	// values that are not read-only (like slices or maps or structs) cannot be
	// constants
	// we can declare them as top-level variables
	b = []byte{3, 4} // this could not be a constant
)

// :show end

func main() {
	// do nothing
}
