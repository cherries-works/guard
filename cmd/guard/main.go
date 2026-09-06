package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/cherries-works/guard/internal/analyzer"
)

var SPINNER = []string{"|", "/", "-", "\\"}

func main() {
	verbose := flag.Bool("verbose", false, "list every advisory affecting each vulnerable dependency")
	flag.BoolVar(verbose, "v", false, "shorthand for --verbose")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: guard [-v] <path>")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		return
	}

	pwd := flag.Arg(0)

	done := make(chan analyzer.Analysis)

	go func() {
		done <- analyzer.Analyzer(pwd)
	}()

	for i := 0; ; i++ {
		select {
		case analysis := <-done:
			fmt.Print("\r\033[K")
			analyzer.PrintAnalysis(analysis, *verbose)
			return

		default:
			fmt.Printf("\rScanning dependencies %s", SPINNER[i%len(SPINNER)])
			time.Sleep(100 * time.Millisecond)
		}
	}
}
