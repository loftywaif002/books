// no output
package main

func main() {
	// :show start
	// you can cast between numbers i.e. integers of various sizes and floating point numbers
	var i1 int32 = 3
	var i2 int = int(i1) // we must explicitly cast from int32 to int
	var f float64 = float64(i1)

	s := "string"
	// we can cast between string and []byte and vice-versa
	// note that unless optimizted by the compiler, this involves allocation
	var d []byte = []byte(s)
	// :show end

	_, _, _ = i2, f, d // silence compiler error about unused variables
}
