package common

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// GzippedReadCloser is a io.ReadCloser for a gzip file
type GzippedReadCloser struct {
	f *os.File
	r io.Reader
}

// Close closes a reader
func (rc *GzippedReadCloser) Close() error {
	return rc.f.Close()
}

// Read reads data from a reader
func (rc *GzippedReadCloser) Read(d []byte) (int, error) {
	return rc.r.Read(d)
}

// OpenGzipped returns io.ReadCloser for a gzip file
func OpenGzipped(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	rc := &GzippedReadCloser{
		f: f,
		r: r,
	}
	return rc, nil
}

// ReadGzipped is like ioutil.ReadFile() but for gzipped files
func ReadGzipped(path string) ([]byte, error) {
	rc, err := OpenGzipped(path)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	d, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// JSONDecodeGzipped reads gzipped file and json-decodes its conent
func JSONDecodeGzipped(path string, v interface{}) error {
	r, err := OpenGzipped(path)
	if err != nil {
		return err
	}
	defer r.Close()
	dec := json.NewDecoder(r)
	return dec.Decode(v)
}

// ReadFileNormalized reads the file and normalized newlines
func ReadFileNormalized(path string) ([]byte, error) {
	d, err := ioutil.ReadFile(path)
	return NormalizeNewlines(d), err
}

// ReadFileAsLines reads a file as lines
func ReadFileAsLines(path string) ([]string, error) {
	d, err := ReadFileNormalized(path)
	if err != nil {
		return nil, err
	}
	s := string(d)
	res := strings.Split(s, "\n")
	return res, nil
}

// GetDirs returns all sub-directories in a dir
func GetDirs(dir string) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, fi := range fileInfos {
		if fi.IsDir() {
			res = append(res, fi.Name())
		}
	}
	return res, nil
}
