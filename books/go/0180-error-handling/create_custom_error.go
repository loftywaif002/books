package main

import (
	"fmt"
)

// :show start

// MyError is a custom error type
type MyError struct {
	msg string
}

func (e *MyError) Error() string {
	return e.msg
}

// :show end

func main() {
	// :show start
	var err error = &MyError{msg: "This is custom error type"}
	fmt.Printf("err: %s\n", err)
	// :show end
}
