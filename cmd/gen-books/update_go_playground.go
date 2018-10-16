package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/kjk/u"

	"github.com/essentialbooks/books/pkg/common"
)

// Sha1ToGoPlaygroundCache maintains sha1 of content to go playground id cache
type Sha1ToGoPlaygroundCache struct {
	cachePath string
	sha1ToID  map[string]string
	nUpdates  int
}

// appends a line to a file
func appendToFile(path string, s string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	_, err = f.WriteString(s)
	if err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func readSha1ToGoPlaygroundCache(path string) *Sha1ToGoPlaygroundCache {
	res := &Sha1ToGoPlaygroundCache{
		cachePath: path,
		sha1ToID:  map[string]string{},
	}
	lines, err := common.ReadFileAsLines(path)
	if err != nil {
		if os.IsNotExist(err) {
			// early detection of "can't create a file" condition
			f, err := os.Create(path)
			panicIfErr(err)
			f.Close()
		}
	}
	for i, s := range lines {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		parts := strings.Split(s, " ")
		panicIf(len(parts) != 2, "unexpected line '%s'", lines[i])
		sha1 := parts[0]
		id := parts[1]
		res.sha1ToID[sha1] = id
	}
	fmt.Printf("Loaded '%s' with %d entries\n", path, len(res.sha1ToID))
	return res
}

// submit the data to Go playground and get share id
func getGoPlaygroundShareID(d []byte) (string, error) {
	uri := "https://play.golang.org/share"
	r := bytes.NewBuffer(d)
	resp, err := http.Post(uri, "text/plain", r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("http.Post returned error code '%s'", err)
	}
	d, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(d)), nil
}

func testGetGoPlaygroundShareIDAndExit() {
	path := "books/go/0230-mutex/rwlock.go"
	d, err := common.ReadFileNormalized(path)
	panicIfErr(err)
	shareID, err := getGoPlaygroundShareID(d)
	panicIfErr(err)
	fmt.Printf("share id: '%s'\n", shareID)
	os.Exit(0)
}

// GetPlaygroundID gets go playground id from content
func (c *Sha1ToGoPlaygroundCache) GetPlaygroundID(d []byte) (string, error) {
	sha1 := u.Sha1HexOfBytes(d)
	id, ok := c.sha1ToID[sha1]
	if ok {
		return id, nil
	}
	id, err := getGoPlaygroundShareID(d)
	if err != nil {
		return "", err
	}
	s := fmt.Sprintf("%s %s\n", sha1, id)
	err = appendToFile(c.cachePath, s)
	if err != nil {
		return "", err
	}
	c.nUpdates++
	return id, nil
}

func getSha1ToGoPlaygroundIDCached(b *Book, d []byte) (string, error) {
	nUpdates := b.sha1ToGoPlaygroundCache.nUpdates
	id, err := b.sha1ToGoPlaygroundCache.GetPlaygroundID(d)
	if err == nil && nUpdates != b.sha1ToGoPlaygroundCache.nUpdates {
		sha1 := u.Sha1HexOfBytes(d)
		fmt.Printf("getSha1ToGoPlaygroundIDCached: %s => %s\n", sha1, id)
	}
	return id, err
}
