// no playground
package main

import (
	"fmt"
	"log"
	"os/exec"
	"sync/atomic"
	"time"
)

func main() {
	// :show start
	cmd := exec.Command("go", "version")

	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	var timedOut int32
	timeout := 1 * time.Millisecond
	stopTimer := time.AfterFunc(timeout, func() {
		cmd.Process.Kill()
		atomic.StoreInt32(&timedOut, 1)
	})

	err = cmd.Wait()
	stopTimer.Stop()
	didTimeout := atomic.LoadInt32(&timedOut) != 0
	fmt.Printf("didTimeout: %v, err: %v\n", didTimeout, err)
	// :show end
}
