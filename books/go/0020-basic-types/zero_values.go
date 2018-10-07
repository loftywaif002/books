package main

import "fmt"

func main() {
	// :show start
	var zeroBool bool
	var zeroInt int
	var zeroF32 float32
	var zeroStr string
	var zeroPtr *int
	var zeroSlice []uint32
	var zeroMap map[string]int
	var zeroInterface interface{}
	var zeroChan chan bool
	var zeroArray [5]int
	type struc struct {
		a int
		b string
	}
	var zeroStruct struc
	var zeroFunc func(bool)

	fmt.Printf("zero bool:       %v\n", zeroBool)
	fmt.Printf("zero int:        %v\n", zeroInt)
	fmt.Printf("zero f32:        %v\n", zeroF32)
	fmt.Printf("zero string:     %#v\n", zeroStr)
	fmt.Printf("zero ptr:        %v\n", zeroPtr)
	fmt.Printf("zero slice:      %v\n", zeroSlice)
	fmt.Printf("zero map:        %#v\n", zeroMap)
	fmt.Printf("zero interface:  %v\n", zeroInterface)
	fmt.Printf("zero channel:    %v\n", zeroChan)
	fmt.Printf("zero array:      %v\n", zeroArray)
	fmt.Printf("zero struct:     %#v\n", zeroStruct)
	fmt.Printf("zero function:   %v\n", zeroFunc)
	// :show end
}
