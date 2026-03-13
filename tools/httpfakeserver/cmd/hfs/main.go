package main

import (
	"os"

	"github.com/suzuito/sandbox2-common-go/tools/httpfakeserver"
)

func main() {
	os.Exit(httpfakeserver.Main())
}
