// allow error
package main

import "fmt"

func main() {
	// :show start
	var m map[string]string

	// you can read read from un-initialized map
	fmt.Printf(`m["foo"] = %s`+"\n", m["foo"])
	_, ok := m["foo"]
	fmt.Printf("ok: %v\n", ok)

	// writing to uninitialized map panics
	m["foo"] = "bar"
	// :show end
}
