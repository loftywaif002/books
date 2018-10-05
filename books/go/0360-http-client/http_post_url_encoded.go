// no playground
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	// :show start
	client := &http.Client{}
	client.Timeout = time.Second * 15

	uri := "https://httpbin.org/post"
	data := url.Values{
		"name":  []string{"John"},
		"email": []string{"john@gmail.com"},
	}
	resp, err := client.PostForm(uri, data)
	if err != nil {
		log.Fatalf("client.PosFormt() failed with '%s'\n", err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll() failed with '%s'\n", err)
	}

	fmt.Printf("PostForm() sent '%s'. Response status code: %d\n", data.Encode(), resp.StatusCode)
	// :show end
}
