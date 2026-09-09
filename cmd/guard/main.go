package main

import (
	"flag"

	"github.com/cherries-works/guard/internal/commands"
	"github.com/cherries-works/guard/internal/utils"
)

func main() {
	scanCommand := flag.NewFlagSet("scan", flag.ExitOnError)
	scanPwd := scanCommand.String("path", ".", "Path to iterate from.")
	scanCommand.StringVar(scanPwd, "p", ".", "Shorthand for --path")
	scanVerbose := scanCommand.Bool("verbose", false, "List every advisory affecting each vulnerable dependency.")
	scanCommand.BoolVar(scanVerbose, "v", false, "Shorthand for --verbose")

	flag.Usage = func() {
		utils.Help()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		return
	}

	cmd := flag.Arg(0)
	switch cmd {
	case "scan":
		commands.Scanner(*scanPwd, *scanVerbose)
	case "help":
		flag.Usage()
	}
}
