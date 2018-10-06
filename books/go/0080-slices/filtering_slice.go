package main

import "fmt"

func filterEvenValuesInPlace(a []int) []int {
	// create a zero-length slice with the same underlying array
	res := a[:0]

	for _, v := range a {
		if v%2 == 0 {
			// collect only wanted values
			res = append(res, v)
		}
	}
	return res
}

func main() {
	// :show start
	a := []int{1, 2, 3, 4}
	res := filterEvenValuesInPlace(a)
	fmt.Printf("%#v\n", res)
	// :show end
}
