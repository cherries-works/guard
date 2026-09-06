package utils

import (
	"fmt"
	"os"
)

const RED = "\033[31m"
const BOLD = "\033[1m"
const DIM = "\033[2m"
const RESET = "\033[0m"

func Title() {
	fmt.Fprintf(os.Stdout, "%s%sCherries Guard%s ───────────────────────────────────── v0.1.0 ────\n", RED, BOLD, RESET)
}

func Help() {
	Title()
	fmt.Fprintf(os.Stdout, "Usage: guard [-v] [--verbose] <path>\n")
	fmt.Fprintf(os.Stdout, "  %-21s List every advisotry affecting each vulnerable dependency.\n", "-v, --verbose")
	fmt.Fprintf(os.Stdout, "  %-21s Prints this.\n", "-h, --help")
	fmt.Fprintf(os.Stdout, "\n")
}
