// no playground
package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// :show start
func writeCSV() error {
	f, err := os.Create("stocks_tmp.csv")
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	records := [][]string{
		{"date", "price", "name"},
		{"2013-02-08", "15,07", "GOOG"},
		{"2013-02-09", "15,09", "GOOG"},
	}
	for _, rec := range records {
		err = w.Write(rec)
		if err != nil {
			f.Close()
			return err
		}
	}

	// csv.Writer might buffer writes for performance so we must
	// Flush to ensure all data has been written to underlying
	// writer
	w.Flush()

	// Flush doesn't return an error. If it failed to write, we
	// can get the error with Error()
	err = w.Error()
	if err != nil {
		return err
	}
	// Close might also fail due to flushing out buffered writes
	err = f.Close()
	return err
}

// :show end

func main() {
	writeCSV()
	d, err := ioutil.ReadFile("stocks_tmp.csv")
	if err != nil {
		log.Fatalf("ioutil.ReadAll() failed with '%s'\n", err)
	}
	fmt.Printf("%s\n", string(d))
	os.Remove("stocks_tmp.csv")
}
