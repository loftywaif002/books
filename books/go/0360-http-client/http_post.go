// no playground
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	// :show start
	client := &http.Client{}
	client.Timeout = time.Second * 15

	uri := "https://httpbin.org/post"
	body := bytes.NewBufferString("text we send")
	resp, err := client.Post(uri, "text/plain", body)
	if err != nil {
		log.Fatalf("client.Post() failed with '%s'\n", err)
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("http.Get() failed with '%s'\n", err)
	}
	fmt.Printf("http.Post() returned statuts code %d, truncated text:\n%s...\n", resp.StatusCode, string(d)[:93])
	// :show end
}
