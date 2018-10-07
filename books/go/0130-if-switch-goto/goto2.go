// allow error
package main

import "fmt"

func main() {
	// :show start
	goto end
	a := 3
	fmt.Printf("a: %d\n", a)
end:
	// :show end
}
