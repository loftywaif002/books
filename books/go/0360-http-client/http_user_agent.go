// no playground
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	// :show start
	userAgent := "Special Agent v 1.0.16"

	client := &http.Client{}
	client.Timeout = time.Second * 15

	uri := "https://httpbin.org/user-agent"
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Fatalf("http.NewRequest() failed with '%s'\n", err)
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("client.Do() failed with '%s'\n", err)
	}

	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll() failed with '%s'\n", err)
	}

	fmt.Printf("Response status code: %d, text:\n%s\n", resp.StatusCode, string(d))
	// :show end
}
