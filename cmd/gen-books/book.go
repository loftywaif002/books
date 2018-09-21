package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/kjk/u"
)

func updateBookAppJS(book *Book) {
	srcName := fmt.Sprintf("app-%s.js", book.titleSafe)
	path := filepath.Join("tmpl", "app.js")
	d, err := ioutil.ReadFile(path)
	maybePanicIfErr(err)
	if err != nil {
		return
	}
	if doMinify {
		d2, err := minifier.Bytes("text/javascript", d)
		maybePanicIfErr(err)
		if err == nil {
			fmt.Printf("Minified %s from %d => %d (saved %d)\n", srcName, len(d), len(d2), len(d)-len(d2))
			d = d2
		}
	}

	d = append(book.tocData, d...)
	sha1Hex := u.Sha1HexOfBytes(d)
	name := nameToSha1Name(srcName, sha1Hex)
	dst := filepath.Join("www", "s", name)
	err = ioutil.WriteFile(dst, d, 0644)
	maybePanicIfErr(err)
	if err != nil {
		return
	}
	book.AppJSURL = "/s/" + name
	fmt.Printf("Created %s\n", dst)
}

var didPrint = false

func printKnownURLS(a []string) {
	if didPrint {
		return
	}
	didPrint = true
	fmt.Printf("%d known urls\n", len(a))
	for _, s := range a {
		fmt.Printf("%s\n", s)
	}
}

// turn partial url like "20381" into a full url like "20381-installing"
func (b *Book) fixupURL(uri string) string {
	// skip uris that are not article/chapter uris
	if strings.Contains(uri, "/") {
		return uri
	}
	for _, known := range b.knownUrls {
		if uri == known {
			return uri
		}
		if strings.HasPrefix(known, uri) {
			//fmt.Printf("fixupURL: %s => %s\n", uri, known)
			return known
		}
	}
	fmt.Printf("fixupURL: didn't fix up: %s\n", uri)
	//printKnownURLS(knownURLS)
	return uri
}

func (b *Book) makeFixupURL() func(uri string) string {
	return func(uri string) string {
		return b.fixupURL(uri)
	}
}
