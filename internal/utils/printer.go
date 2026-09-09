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
	fmt.Fprintf(os.Stdout, " > %-20s %-20s\n", "scan", "Scans your project for vulnerable/outdated dependencies.")
	fmt.Fprintf(os.Stdout, "     %s%-20s %-20s%s\n", DIM, "--path (-p) [string]", "Specify which path to check the dependencies from (default \".\").", RESET)
	fmt.Fprintf(os.Stdout, "     %s%-20s %-20s%s\n", DIM, "--verbose (-v)", "Makes output more verbose.", RESET)
	fmt.Fprintf(os.Stdout, " > %-20s %-20s\n", "help", "Prints this.")
	fmt.Fprintf(os.Stdout, "\n")
}
