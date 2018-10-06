// no output
package main

func main() {
	// :show start
	src := make(map[string]int)
	src["one"] = 1
	src["two"] = 2

	dst := make(map[string]int)

	for key, value := range src {
		dst[key] = value
	}

	// :show end
}
