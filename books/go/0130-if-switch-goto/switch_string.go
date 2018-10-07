package main

import "fmt"

func main() {
	// :show start
	s := "foo"
	switch s {
	case "foo":
		fmt.Printf("s is 'foo'\n")
	case "bar":
		fmt.Printf("s is 'bar'\n")
	}
	// :show end
}
