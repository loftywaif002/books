// allow error
package main

func main() {
	// :show start
	ch := make(chan int)
	close(ch)
	ch <- 5 // panics
	// :show end
}
