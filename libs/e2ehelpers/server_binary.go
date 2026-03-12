package e2ehelpers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type RunServerInput struct {
	Args []string
	Envs []string
}

func RunServer(
	ctx context.Context,
	filePathBin string,
	input *RunServerInput,
	healthCheckFunc func() error,
) (func() (exitCode int, stdout string, stderr string, err error), bool) {
	if filePathBin == "" {
		panic("filePathBin is empty")
	}

	cmd := exec.CommandContext(
		ctx,
		filePathBin,
		input.Args...,
	)

	cmd.Env = append(
		os.Environ(),
		input.Envs...,
	)

	stdout, stderr := bytes.NewBufferString(""), bytes.NewBufferString("")
	printStdoutStderr := func() {
		fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
		fmt.Println("@@@@@@@ STDOUT @@@@@@@")
		fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
		fmt.Println(stdout.String())
		fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
		fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
		fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
		fmt.Println()

		if stderr.Len() > 0 {
			fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
			fmt.Println("@@@@@@@ STDERR @@@@@@@")
			fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
			fmt.Println(stderr.String())
			fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
			fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
			fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
		}
	}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Start()
	if err != nil {
		var exiterr *exec.ExitError
		if !errors.As(err, &exiterr) {
			panic(fmt.Sprintf("%s %s: %s", filePathBin, strings.Join(input.Args, " "), err.Error()))
		}
	}

	shutdown := func() (int, string, string, error) {
		defer func() {
			printStdoutStderr()
		}()

		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			return 0, "", "", err
		}

		if err := cmd.Wait(); err != nil {
			var exiterr *exec.ExitError
			if !errors.As(err, &exiterr) {
				return 0, "", "", err
			}
		}

		return cmd.ProcessState.ExitCode(), stdout.String(), stderr.String(), nil
	}

	if err := healthCheckFunc(); err != nil {
		fmt.Fprintf(os.Stderr, "health check error: %s\n", err.Error())
		return shutdown, false
	}

	return shutdown, true
}

func CheckHTTPServerHealth(
	ctx context.Context,
	u string,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	cli := http.DefaultClient

	for {
		time.Sleep(time.Millisecond * 500)
		select {
		case <-ctx.Done():
			return errors.New("health check is failed")
		default:
			fmt.Printf("health check GET %s\n", u)

			res, err := cli.Get(u)
			if err != nil {
				fmt.Printf("failed to health check: %v\n", err)
				continue
			}

			res.Body.Close() //nolint:errcheck
			if res.StatusCode == http.StatusOK {
				return nil
			}

			fmt.Printf("failed to health check with http error: %d\n", res.StatusCode)
		}
	}
}
