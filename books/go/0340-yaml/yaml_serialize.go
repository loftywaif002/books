// no playground
package main

import (
	"fmt"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// :show start
type Person struct {
	fullName string
	Name     string
	Age      int    `yaml:"age"`
	City     string `yaml:"city"`
}

// :show end

func main() {
	// :show start
	p := Person{
		Name: "John",
		Age:  37,
		City: "SF",
	}
	d, err := yaml.Marshal(&p)
	if err != nil {
		log.Fatalf("yaml.Marshal failed with '%s'\n", err)
	}
	fmt.Printf("Person in YAML:\n%s\n", string(d))
	// :show end
}
