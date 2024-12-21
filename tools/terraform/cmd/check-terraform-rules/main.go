package main

import (
	"flag"
	"fmt"
	"os"
)

var usageString = `This command verifies whether terraform files is preventing rules. If violation is occured, exit not zero.
`

func usage() {
	fmt.Fprintln(os.Stderr, usageString)
	flag.PrintDefaults()
}

func main() {
	var dirPathBase string

	flag.StringVar(&dirPathBase, "d", "", "base directory path")

	flag.Parse()

	if dirPathBase == "" {
		fmt.Fprint(os.Stderr, "-d is required")
		os.Exit(1)
	}
}
