package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/suzuito/sandbox2-common-go/tools/fakecmd/domains"
)

type code uint8

const (
	codeFakeCMDError code = 127
)

var logger *log.Logger

func init() {
	logger = newErrorLogger()
}

func newErrorLogger() *log.Logger {
	prefix := "FAKE_CMD_ERROR "
	var w io.Writer = os.Stderr

	s := strings.ToLower(os.Getenv("FAKECMD_ERROR_LOG"))
	switch s {
	case "discard":
		w = io.Discard
	}

	return log.New(w, prefix, 0)
}

func getLock(
	filePathProcessing string,
) (func(), error) {
	var err error
	var fileProcessing *os.File

	if _, err := os.Stat(filePathProcessing); err == nil {
		return nil, fmt.Errorf("processing file already exists")
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to os.Stat: %w", err)
	}

	fileProcessing, err = os.Create(filePathProcessing)
	if err != nil {
		return nil, fmt.Errorf("failed to os.Create: %s: %w", filePathProcessing, err)
	}
	defer fileProcessing.Close()

	return func() {
		os.Remove(filePathProcessing)
	}, nil
}

func main1() code {

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	dirPath := domains.DirPathFakeCommand(filepath.Dir(os.Args[0]))

	behaviors := domains.Behaviors{}
	behaviorsBytes, err := os.ReadFile(dirPath.FilePathBehaviors())
	if err != nil {
		logger.Printf("failed to read behavior file: %s: %s\n", dirPath.FilePathBehaviors(), err)
		return codeFakeCMDError
	}
	if err := json.Unmarshal(behaviorsBytes, &behaviors); err != nil {
		logger.Printf("failed to unmarshal behavior file: %s: %s\n", dirPath.FilePathBehaviors(), err)
		return codeFakeCMDError
	}
	if len(behaviors) <= 0 {
		logger.Println("no behaviors")
		return codeFakeCMDError
	}

	releaseLock, err := getLock(dirPath.FilePathProcessing())
	if err != nil {
		logger.Printf("failed to get lock: %s: %s\n", dirPath.FilePathProcessing(), err)
		return codeFakeCMDError
	}
	defer releaseLock()

	state := domains.State{}
	contentState, err := os.ReadFile(dirPath.FilePathState())
	if os.IsNotExist(err) {
	} else if err != nil {
		logger.Printf("failed to read state file: %s: %s\n", dirPath.FilePathState(), err)
		return codeFakeCMDError
	} else {
		if err := json.Unmarshal(contentState, &state); err != nil {
			logger.Printf("failed to unmarshal state file: %s: %s\n", dirPath.FilePathState(), err)
			return codeFakeCMDError
		}
	}

	closeState := func() {
		state.ExecutedHistories = append(state.ExecutedHistories, domains.ExecutedHistory{})
		b, err := json.Marshal(state)
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile(dirPath.FilePathState(), b, 0755); err != nil {
			panic(err)
		}
	}

	if len(behaviors) <= len(state.ExecutedHistories) {
		// undefined executions over state.TimesExecuted
		logger.Printf("all expected executions are done: expected=%d histories=%d\n", len(behaviors), len(state.ExecutedHistories))
		return codeFakeCMDError
	}
	behavior := behaviors[len(state.ExecutedHistories)]
	switch behavior.Type {
	case domains.BehaviorTypeStdoutStderrExitCode:
		if behavior.BehaviorStdoutStderrExitCode == nil {
			// type is BehaviorTypeStdoutStderrExitCode but nil
			return codeFakeCMDError
		}
		defer closeState()
		fmt.Fprintf(os.Stdout, behavior.BehaviorStdoutStderrExitCode.Stdout)
		fmt.Fprintf(os.Stderr, behavior.BehaviorStdoutStderrExitCode.Stderr)
		return code(behavior.BehaviorStdoutStderrExitCode.ExitCode)
	default:
		// unknown type
		return codeFakeCMDError
	}
}

func main() {
	os.Exit(int(main1()))
}
