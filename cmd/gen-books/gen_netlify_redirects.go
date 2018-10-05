package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	// https://www.netlify.com/docs/headers-and-basic-auth/#custom-headers
	netlifyHeaders = `
# long-lived caching
/s/*
  Cache-Control: max-age=31536000
/*
  X-Content-Type-Options: nosniff
  X-Frame-Options: DENY
  X-XSS-Protection: 1; mode=block
`
)

func genNetlifyHeaders() {
	path := filepath.Join("www", "_headers")
	err := ioutil.WriteFile(path, []byte(netlifyHeaders), 0644)
	panicIfErr(err)
}

func genNetlifyRedirectsForBook(b *Book) []string {
	var res []string

	pages := b.GetAllPages()
	for _, page := range pages {
		id := page.NotionID
		uri := page.URLLastPath()
		s := fmt.Sprintf(`/essential/%s/%s* /essential/%s/%s 302`, b.Dir, id, b.Dir, uri)
		// TODO: also add redirects for old article ids (only needed for Go book)
		// alternatively just forget about it
		res = append(res, s)
	}

	// catch-all redirect for all other missing pages
	s := fmt.Sprintf(`/essential/%s/* /essential/%s/404.html 404`, b.Dir, b.Dir)
	res = append(res, s)
	res = append(res, "")
	return res
}

func genNetlifyRedirects() {
	var a []string
	for _, b := range books {
		ab := genNetlifyRedirectsForBook(b)
		a = append(a, ab...)
	}
	s := strings.Join(a, "\n")
	path := filepath.Join("www", "_redirects")
	err := ioutil.WriteFile(path, []byte(s), 0644)
	panicIfErr(err)
}
