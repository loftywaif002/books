// allow error
package main

import "io"

// :show start
type MyReadCloser struct {
}

func (rc *MyReadCloser) Read(d []byte) (int, error) {
	return 0, nil
}

var _ io.ReadCloser = &MyReadCloser{}

// :show end

func main() {
}
