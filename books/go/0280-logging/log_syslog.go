// no output
package main

import (
	"log"
	"log/syslog"
)

func main() {
	// :show start
	syslogger, err := syslog.New(syslog.LOG_INFO, "syslog_example")
	if err != nil {
		log.Fatalln(err)
	}

	log.SetOutput(syslogger)
	log.Println("Log entry")
	// :show end
}
