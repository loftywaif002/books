// no playground
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// :show start
type postItem struct {
	Score int    `json:"score"`
	Link  string `json:"link"`
}

type postsType struct {
	Items []postItem `json:"items"`
}

// :show end
func main() {
	// :show start
	values := url.Values{
		"order": []string{"desc"},
		"sort":  []string{"activity"},
		"site":  []string{"stackoverflow"},
	}

	// URL parameters can also be programmatically set
	values.Set("page", "1")
	values.Set("pagesize", "10")

	uri := "https://api.stackexchange.com/2.2/posts?"
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Get(uri + values.Encode())
	if err != nil {
		log.Fatalf("http.Get() failed with '%s'\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		d, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Request was '%s' (%d) and not OK (200). Body:\n%s\n", resp.Status, resp.StatusCode, string(d))
	}

	dec := json.NewDecoder(resp.Body)
	var p postsType
	err = dec.Decode(&p)
	if err != nil {
		log.Fatalf("dec.Decode() failed with '%s'\n", err)
	}

	fmt.Println("Top 10 most recently active StackOverflow posts:")
	fmt.Println("Score", "Link")
	for _, post := range p.Items {
		fmt.Println(post.Score, post.Link)
	}
	// :show end
}
