package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/cherries-works/guard/internal/analyzer"
	"github.com/cherries-works/guard/internal/utils"
)

var SPINNER = []string{"|", "/", "-", "\\"}

func main() {
	verbose := flag.Bool("verbose", false, "List every advisory affecting each vulnerable dependency.")
	flag.BoolVar(verbose, "v", false, "Shorthand for --verbose")

	help := flag.Bool("help", false, "Prints help.")
	flag.BoolVar(help, "h", false, "Shorthand for --help")

	flag.Usage = func() {
		utils.Help()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		return
	}

	if *help {
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
