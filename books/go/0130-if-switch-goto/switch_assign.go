package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// without seeding rand.Intn will always return the same sequence
	rand.Seed(time.Now().UnixNano())
	// :show start
	switch n := rand.Intn(9); n {
	case 1, 2, 3:
		fmt.Printf("case 1, 2, 3: n is %d\n", n)
	case 4, 5:
		fmt.Printf("case 4, 5: n is %d\n", n)
	default:
		fmt.Printf("default: n is %d\n", n)
	}
	// :show end
}
