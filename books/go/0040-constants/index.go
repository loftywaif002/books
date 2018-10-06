// no output
package main

// :show start
// Greeting is an exported (public) string constant
const Greeting string = "Hello World"

// we can group const declarations
const (
	// years is an unexported (package private) int constant
	years int  = 10
	truth bool = true
)

// :show end

func main() {
	// do  nothing
}
