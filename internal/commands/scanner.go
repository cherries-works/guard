package commands

import (
	"fmt"
	"time"

	"github.com/cherries-works/guard/internal/analyzer"
)

var SPINNER = []string{"|", "/", "-", "\\"}

func Scanner(pwd string, verbose bool) {
	done := make(chan analyzer.Analysis)
	go func() {
		done <- analyzer.Analyzer(pwd)
	}()

	for i := 0; ; i++ {
		select {
		case analysis := <-done:
			fmt.Print("\r\033[K")
			analyzer.PrintAnalysis(analysis, verbose)
			return

		default:
			fmt.Printf("\rScanning dependencies %s", SPINNER[i%len(SPINNER)])
			time.Sleep(100 * time.Millisecond)
		}
	}
}
