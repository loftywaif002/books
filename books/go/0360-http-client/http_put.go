// no playground
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type User struct {
	Name  string
	Email string
}

func main() {
	// :show start
	user := User{
		Name:  "John Doe",
		Email: "johndoe@example.com",
	}

	d, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("json.Marshal() failed with '%s'\n", err)
	}

	client := &http.Client{}
	client.Timeout = time.Second * 15

	uri := "https://httpbin.org/put"
	body := bytes.NewBuffer(d)
	req, err := http.NewRequest(http.MethodPut, uri, body)
	if err != nil {
		log.Fatalf("http.NewRequest() failed with '%s'\n", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("client.Do() failed with '%s'\n", err)
	}

	defer resp.Body.Close()
	d, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll() failed with '%s'\n", err)
	}

	fmt.Printf("Response status code: %d, text:\n%s\n", resp.StatusCode, string(d))
	// :show end
}
