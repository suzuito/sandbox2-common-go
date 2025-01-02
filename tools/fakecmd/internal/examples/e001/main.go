package main

import (
	"fmt"
	"os"

	"github.com/suzuito/sandbox2-common-go/tools/fakecmd/domains"
)

func main() {
	filePathFakeCMD := fmt.Sprintf("%s/bin/main", os.Getenv("GOPATH"))
	if _, err := os.Stat(filePathFakeCMD); err != nil {
		panic(err)
	}

	faker := domains.NewExternalCommandFaker(filePathFakeCMD)
	fcmd, err := faker.Add("/tmp/hoge", domains.Behaviors{
		{
			Type: domains.BehaviorTypeStdoutStderrExitCode,
			BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
				Stdout:   "hello world!",
				ExitCode: 10,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	if err := fcmd.Init(true); err != nil {
		panic(err)
	}

	fmt.Println(fcmd.DirPath())
}
