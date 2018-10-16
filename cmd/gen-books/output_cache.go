package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/essentialbooks/books/pkg/kvstore"
	"github.com/google/shlex"
	"github.com/kjk/u"
)

const maxOutputFileSize = 1024 * 128 // 128 kB

type cachedOutputFile struct {
	path string
	doc  kvstore.Doc
	no   int
}

func getCurrentOutputCacheFile(b *Book) *cachedOutputFile {
	n := len(b.cachedOutputFiles) - 1
	if n >= 0 {
		cof := b.cachedOutputFiles[n]
		if getDocSize(cof.doc) < maxOutputFileSize {
			return cof
		}
	}
	fileNo := len(b.cachedOutputFiles) + 1
	name := fmt.Sprintf("cached_output_%d.txt", fileNo)
	path := filepath.Join(b.OutputCacheDir(), name)
	cof := &cachedOutputFile{
		path: path,
		doc:  nil,
		no:   fileNo,
	}
	b.cachedOutputFiles = append(b.cachedOutputFiles, cof)
	fmt.Printf("Created new cachedOutputFile. path: '%s'\n", path)
	return cof
}

func getDocSize(doc kvstore.Doc) int {
	size := 0
	for _, kv := range doc {
		size += len(kv.Key)
		size += len(kv.Value)
	}
	return size
}

func isCachedOutputFile(path string) bool {
	return strings.Contains(path, "cached_output_") && strings.HasSuffix(path, ".txt")
}

// given cached_output_${no}.txt return ${no}
func cachedFileNo(path string) int {
	parts := strings.Split(path, "_")
	s := parts[len(parts)-1]
	// now is ${no}.txt
	parts = strings.Split(s, ".")
	n, err := strconv.Atoi(parts[0])
	u.PanicIfErr(err)
	return n
}

// files are cached_output_${no}.txt
func reloadCachedOutputFilesMust(b *Book) {
	b.sha1ToCachedOutputFile = make(map[string]*cachedOutputFile)

	fileInfos, err := ioutil.ReadDir(b.OutputCacheDir())
	u.PanicIfErr(err)
	for _, fi := range fileInfos {
		if fi.IsDir() {
			continue
		}
		if fi.Name() == "sha1_to_go_playground_id.txt" {
			continue
		}
		if !isCachedOutputFile(fi.Name()) {
			u.PanicIf(true, "'%s' is not a file with cached output", fi.Name())
			continue
		}
		path := filepath.Join(b.OutputCacheDir(), fi.Name())
		doc, err := kvstore.ParseKVFile(path)
		u.PanicIfErr(err)
		f := &cachedOutputFile{
			path: path,
			doc:  doc,
			no:   cachedFileNo(path),
		}
		b.cachedOutputFiles = append(b.cachedOutputFiles, f)
	}
	fmt.Printf("loaded %d cached output files\n", len(b.cachedOutputFiles))
	if len(b.cachedOutputFiles) == 0 {
		return
	}
	sort.Slice(b.cachedOutputFiles, func(i, j int) bool {
		n1 := b.cachedOutputFiles[i].no
		n2 := b.cachedOutputFiles[j].no
		return n1 < n2
	})
	//fmt.Printf("%#v\n", cachedOutputFiles)
	for _, cfo := range b.cachedOutputFiles {
		for _, kv := range cfo.doc {
			sha1 := kv.Key
			b.sha1ToCachedOutputFile[sha1] = cfo
		}
	}
	fmt.Printf("%d cached files\n", len(b.sha1ToCachedOutputFile))
}

func findOutputBySha1(cof *cachedOutputFile, sha1Hex string) string {
	for _, kv := range cof.doc {
		if sha1Hex == kv.Key {
			return kv.Value
		}
	}
	panicIf(true, "didn't find '%s' in '%s'\n", sha1Hex, cof.path)
	return ""
}

func saveCachedOutputFile(cof *cachedOutputFile) {
	doc := cof.doc
	sort.Slice(doc, func(i, j int) bool {
		k1 := doc[i].Key
		k2 := doc[j].Key
		return k1 < k2
	})
	var recs []string
	for _, kv := range doc {
		s := kvstore.SerializeLong(kv.Key, kv.Value)
		recs = append(recs, s)
	}
	s := strings.Join(recs, "")
	err := ioutil.WriteFile(cof.path, []byte(s), 0644)
	u.PanicIfErr(err)
	fmt.Printf("Wrote '%s'\n", cof.path)
}

func saveCachedOutputFiles(b *Book) {
	for _, cof := range b.cachedOutputFiles {
		saveCachedOutputFile(cof)
	}
	reloadCachedOutputFilesMust(b)
}

// runs `go run ${path}` and returns captured output`
func getGoOutput(path string) (string, error) {
	dir, fileName := filepath.Split(path)
	cmd := exec.Command("go", "run", fileName)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func getRunCmdOutput(path string, runCmd string) (string, error) {
	parts, err := shlex.Split(runCmd)
	maybePanicIfErr(err)
	if err != nil {
		return "", err
	}
	exeName := parts[0]
	parts = parts[1:]
	var parts2 []string
	srcDir, srcFileName := filepath.Split(path)

	// remove empty lines and replace variables
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		switch part {
		case "$file":
			part = srcFileName
		}
		parts2 = append(parts2, part)
	}
	//fmt.Printf("getRunCmdOutput: running '%s' with args '%#v'\n", exeName, parts2)
	cmd := exec.Command(exeName, parts2...)
	cmd.Dir = srcDir
	out, err := cmd.CombinedOutput()
	//fmt.Printf("getRunCmdOutput: out:\n%s\n", string(out))
	return string(out), err
}

func stripCurrentPathFromOutput(s string) string {
	path, err := filepath.Abs(".")
	u.PanicIfErr(err)
	return strings.Replace(s, path, "", -1)
}

// it executes a code file and captures the output
// optional runCmd says
func getOutput(path string, runCmd string) (string, error) {
	if runCmd != "" {
		//fmt.Printf("Found :run cmd '%s' in '%s'\n", runCmd, path)
		s, err := getRunCmdOutput(path, runCmd)
		return stripCurrentPathFromOutput(s), err
	}

	// do default
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".go" {
		s, err := getGoOutput(path)
		return stripCurrentPathFromOutput(s), err
	}
	return "", fmt.Errorf("getOutput(%s): files with extension '%s' are not supported", path, ext)
}

func getOutputCachedForReplit(b *Book, replit *Replit, sf *SourceFile) error {
	if sf.Directive.NoOutput {
		return nil
	}

	sha1Hex := u.Sha1HexOfBytes(sf.Data)

	cfo := b.sha1ToCachedOutputFile[sha1Hex]
	if cfo != nil {
		sf.Output = findOutputBySha1(cfo, sha1Hex)
		return nil
	}

	// save files to temp directory
	dir := os.TempDir()
	defer os.RemoveAll(dir)

	if len(replit.files) == 1 {
		name := replit.files[0].name
		path := filepath.Join(dir, name)

		sf.FileName = name
		sf.Path = path
	}

	for _, rf := range replit.files {
		name := rf.name
		path := filepath.Join(dir, name)
		if sf.Path == "" && strings.HasPrefix(name, "main.") {
			sf.FileName = name
			sf.Path = path
		}
		d := []byte(rf.data)
		err := ioutil.WriteFile(path, d, 0644)
		if err != nil {
			return err
		}
	}

	panicIf(sf.Path == "", "no main file")
	return getOutputCached(b, sf)
}

// for a given file, get output of executing this command
// We cache this as it is the most expensive part of rebuilding books
// If allowError is true, we silence an error from executed command
// This is useful when e.g. executing "go run" on a program that is
// intentionally not valid.
func getOutputCached(b *Book, sf *SourceFile) error {
	if sf.Directive.NoOutput {
		return nil
	}

	sha1Hex := u.Sha1HexOfBytes(sf.Data)

	cfo := b.sha1ToCachedOutputFile[sha1Hex]
	if cfo != nil {
		sf.Output = findOutputBySha1(cfo, sha1Hex)
		return nil
	}

	path := sf.Path
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json", ".csv", ".yml", ".xml":
		return nil
	}

	// fmt.Printf("loadFileCached('%s') failed with '%s'\n", outputPath, err)
	s, err := getOutput(path, sf.RunCmd)
	if err != nil {
		if !sf.Directive.AllowError {
			fmt.Printf("getOutput('%s'), output is:\n%s\n", path, s)
			return err
		}
		err = nil
	}

	fmt.Printf("Got output '%s' for '%s'\n", sha1Hex, path)
	cof := getCurrentOutputCacheFile(b)
	cof.doc = kvstore.ReplaceOrAppend(cof.doc, sha1Hex, s)
	return nil
}
