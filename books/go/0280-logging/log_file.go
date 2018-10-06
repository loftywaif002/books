// no output
package main

import (
	"log"
	"os"
)

func main() {
	// :show start
	logfile, err := os.OpenFile("test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("os.OpenFile() failed with '%s\n", err)
	}
	defer logfile.Close()

	log.SetOutput(logfile)
	log.Println("Log entry")
	// :show end
	os.Remove("test.log")
}
