package main

import (
	"bytes"
	"fmt"
)

// :show start
// NormalizeNewLines normalizes \r\n (windows) and \r (mac)
// into \n (unix)
func NormalizeNewlines(d []byte) []byte {
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)
	return d
}

// :show end

func main() {
	d := []byte("new\r\nline")
	d = NormalizeNewlines(d)
	fmt.Printf("%#v\n", string(d))
}
