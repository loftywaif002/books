package main

import "fmt"

func main() {
	// :show start
	var s string // empty string ""
	s1 := "string\nliteral\nwith\tescape characters\n"
	s2 := `raw string literal
which doesnt't recgonize escape characters like \n
`
	fmt.Printf("sum of strings\n'%s'\n", s+s1+s2)
	// :show end
}
