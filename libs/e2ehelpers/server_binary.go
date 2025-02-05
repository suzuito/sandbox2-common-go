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
) func() error {
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

	if err := healthCheckFunc(); err != nil {
		printStdoutStderr()
		panic(fmt.Errorf("health check error: %w", err))
	}

	return func() error {
		defer func() {
			printStdoutStderr()
		}()

		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			return err
		}

		if err := cmd.Wait(); err != nil {
			return err
		}

		return nil
	}
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
			res, err := cli.Get(u)
			if err != nil {
				continue
			}

			res.Body.Close()
			if res.StatusCode == http.StatusOK {
				return nil
			}
		}
	}
}
