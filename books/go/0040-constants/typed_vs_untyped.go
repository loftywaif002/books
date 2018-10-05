// allow error
package main

func main() {
	// :show start
	const untypedNumber = 345

	// if this was a variable declaration, untypedNumber would have its type
	// inferred as int.
	// since it's a const, it remains untyped until it's e.g. assigned to a variable
	var i int = untypedNumber      // no need to cast to int
	var u16 uint16 = untypedNumber // no need to cast to uint16
	var f float64 = untypedNumber  // no need to cast to float64

	// incompatible assignments are detected by the compiler
	// 345 is too big to fit in int8 and compiler detects that
	var b int8 = untypedNumber

	// :show end

	// silence compiler error about unused variables
	_, _, _, _ = i, u16, f, b
}
