package main

import (
	"flag"
	"fmt"
	"os"
)

var usageString = `This command execute "terraform plan" command for changed file on Github PR.
`

func usage() {
	fmt.Fprintln(os.Stderr, usageString)
	flag.PrintDefaults()
}

func main() {
	var dirPathBase string
	var dirPathOutput string
	var errorIfApplyRequired bool

	flag.StringVar(&dirPathBase, "d", "", "Base directory path")
	flag.StringVar(&dirPathOutput, "o", "", "Directory path in which results is saved. Results include Terraform module paths to be required applying.")
	flag.BoolVar(&errorIfApplyRequired, "e", false, "Whether exit error ")

}
