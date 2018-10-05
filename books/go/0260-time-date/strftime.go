// no playground
package main

import (
	"fmt"
	"time"

	strftime "github.com/jehiah/go-strftime"
)

func main() {
	// :show start
	t := time.Date(2017, 9, 4, 3, 38, 45, 0, time.UTC)
	fmt.Println(strftime.Format("%Y-%m-%d %H:%M:%S", t))
	// :show end
}
