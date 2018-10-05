// allow error, no playground
package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	// :show start
	client := &http.Client{}
	client.Timeout = time.Millisecond * 100

	uri := "https://httpbin.org/delay/3"
	resp, err := client.Get(uri)
	if err != nil {
		log.Fatalf("http.Get() failed with '%s'\n", err)
	}
	// :show end
	resp.Body.Close()
}
