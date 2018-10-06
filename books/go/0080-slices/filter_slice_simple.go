package main

import "fmt"

// :show start
func filterEvenValues(a []int) []int {
	var res []int
	for _, el := range a {
		if el%2 == 0 {
			continue
		}
		res = append(res, el)
	}
	return res
}

// :show end

func main() {
	// :show start
	a := []int{1, 2, 3, 4}
	res := filterEvenValues(a)
	fmt.Printf("%#v\n", res)
	// :show end
}
