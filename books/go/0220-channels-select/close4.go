// allow error
package main

func main() {
	// :show start
	ch := make(chan string)
	close(ch)
	close(ch)
	// :show end
}
