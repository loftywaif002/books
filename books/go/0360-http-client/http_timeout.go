// allow error, no playground
package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	// :show start

	uri := "https://httpbin.org/delay/3"
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalf("http.NewRequest() failed with '%s'\n", err)
	}

	ctx, _ := context.WithTimeout(context.TODO(), time.Millisecond*100)
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Fatalf("http.DefaultClient.Do() failed with:\n'%s'\n", err)
	}
	defer resp.Body.Close()
	// :show end
}
