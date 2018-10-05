// no playground
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// :show start
	uri := "https://httpbin.org/html"
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalf("http.Get() failed with '%s'\n", err)
	}

	// it's important to close resp.Body or else we'll leak network connection
	// it must be done after checking for error because in error case
	// resp.Body can be nil
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll() failed with '%s'\n", err)
	}

	contentType := resp.Header.Get("Content-Type")
	fmt.Printf("http.Get() returned content of type '%s' and size %d bytes.\nStatus code: %d\n", contentType, len(d), resp.StatusCode)

	// getting page that doesn't exist return 404
	uri = "https://httpbin.org/page-doesnt-exist"
	resp, err = http.Get(uri)
	if err != nil {
		log.Fatalf("http.Get() failed with '%s'\n", err)
	}

	contentType = resp.Header.Get("Content-Type")
	fmt.Printf("\nhttp.Get() returned content of type '%s' and size %d bytes.\nStatus code: %d\n", contentType, len(d), resp.StatusCode)

	// acessing non-existent host fails
	uri = "http://website.not-exists.as/index.html"
	resp, err = http.Get(uri)
	if err != nil {
		fmt.Printf("\nhttp.Get() failed with: '%s'\nresp: %v\n", err, resp)
	}
	// :show end
}
