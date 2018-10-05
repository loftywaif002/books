// no playground
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// :show start

// ReadLines reads all lines from a file
func ReadLines(filePath string) ([]string, error) {
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

func main() {
	path := "read_lines.go"
	lines, err := ReadLines(path)
	if err != nil {
		log.Fatalf("ReadLines failed with '%s'\n", err)
	}
	fmt.Printf("File %s has %d lines\n", path, len(lines))
}

// :show end
