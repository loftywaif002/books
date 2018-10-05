// allow error
package main

import (
	"errors"
	"fmt"
)

// :show start
// FmtArgs formats args as a string. First argument should be format string
// and the rest are arguments to the format
func FmtArgs(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	format := args[0].(string)
	if len(args) == 1 {
		return format
	}
	return fmt.Sprintf(format, args[1:]...)
}

func panicWithMsg(defaultMsg string, args ...interface{}) {
	s := FmtArgs(args...)
	if s == "" {
		s = defaultMsg
	}
	fmt.Printf("%s\n", s)
	panic(s)
}

// PanicIf panics if cond is true
func PanicIf(cond bool, args ...interface{}) {
	if !cond {
		return
	}
	panicWithMsg("PanicIf: condition failed", args...)
}

// PanicIfErr panics if err is not nil
func PanicIfErr(err error, args ...interface{}) {
	if err == nil {
		return
	}
	panicWithMsg(err.Error(), args...)
}

func main() {
	PanicIfErr(nil)                              // nothing happens
	PanicIfErr(errors.New("there was an error")) // will panic
}

// :show end
