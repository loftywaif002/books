package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	destDir = "out"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func getDirs() []string {
	fileInfos, err := ioutil.ReadDir(".")
	panicIfErr(err)
	res := []string{}
	for _, fi := range fileInfos {
		if !fi.IsDir() {
			continue
		}
		res = append(res, fi.Name())
	}
	return res
}

func mkdirForFile(filePath string) error {
	dir := filepath.Dir(filePath)
	return os.MkdirAll(dir, 0755)
}

func readLines(filePath string) ([]string, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	res := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Bytes()
		res = append(res, string(line))
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func tweakAndCopyFile(dst string, src string) {
	err := mkdirForFile(dst)
	panicIfErr(err)

	lines, err := readLines(src)
	panicIfErr(err)
	if lines[0] == "---" {
		lines = lines[1:]
	}

	d := strings.Join(lines, "\n")
	err = ioutil.WriteFile(dst, []byte(d), 0644)
	panicIfErr(err)
}

func copyMdFilesFromDir(dir string) {
	fileInfos, err := ioutil.ReadDir(dir)
	panicIfErr(err)
	parts := strings.Split(dir, "-")
	suffix := parts[0]
	for _, fi := range fileInfos {
		if fi.IsDir() {
			continue
		}
		name := fi.Name()
		if !strings.HasSuffix(name, ".md") {
			continue
		}
		src := filepath.Join(dir, name)
		dstName := suffix + "-" + name
		dst := filepath.Join(destDir, dstName)
		tweakAndCopyFile(dst, src)
	}
}

func main() {
	dirs := getDirs()
	for _, dir := range dirs {
		copyMdFilesFromDir(dir)
	}
}
