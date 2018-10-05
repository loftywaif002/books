// no playground
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// :show start
	f, err := os.Open("stocks.csv")
	if err != nil {
		log.Fatalf("os.Open() failed with '%s'\n", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	var record []string
	nRecords := 0
	for {
		record, err = r.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		nRecords++
		if nRecords < 5 {
			fmt.Printf("Record: %#v\n", record)
		}
	}
	if err != nil {
		log.Fatalf("r.Read() failed with '%s'\n", err)
	}
	fmt.Printf("Read %d records\n", nRecords)
	// :show end
}
