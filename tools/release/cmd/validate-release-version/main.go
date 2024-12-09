package main

import (
	"flag"
	"log"
)

func main() {
	var version string
	flag.StringVar(
		&version,
		"version",
		"",
		"released version",
	)

	flag.Parse()

	if version == "" {
		log.Fatal("-version is empty")
	}

}
