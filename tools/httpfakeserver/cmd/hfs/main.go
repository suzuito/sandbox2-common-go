package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/suzuito/sandbox2-common-go/tools/httpfakeserver"
)

func main() {
	port := 8080
	portString := os.Getenv("PORT")
	if portString != "" {
		var err error
		port, err = strconv.Atoi(portString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to convert port into int\n")
			os.Exit(1)
		}
	}

	basePathAdmin := os.Getenv("BASE_PATH_ADMIN")
	if basePathAdmin == "" {
		basePathAdmin = "/admin"
	}

	os.Exit(httpfakeserver.Main(httpfakeserver.Options{
		Port:          port,
		BasePathAdmin: basePathAdmin,
	}))
}
