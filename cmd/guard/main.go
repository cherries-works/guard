package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cherries-works/guard/internal/analyzer"
)

var SPINNER = []string{"|", "/", "-", "\\"}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: guard <path>")
		return
	}

	pwd := os.Args[1]

	done := make(chan analyzer.Analysis)

	go func() {
		done <- analyzer.Analyzer(pwd)
	}()

	for i := 0; ; i++ {
		select {
		case analysis := <-done:
			fmt.Print("\r\033[K")
			analyzer.PrintAnalysis(analysis)
			return

		default:
			fmt.Printf("\rScanning dependencies %s", SPINNER[i%len(SPINNER)])
			time.Sleep(100 * time.Millisecond)
		}
	}
}
