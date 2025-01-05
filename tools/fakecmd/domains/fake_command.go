package domains

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type DirPathFakeCommand string

func (t DirPathFakeCommand) String() string {
	return string(t)
}

func (t DirPathFakeCommand) FilePathCommand() string {
	return fmt.Sprintf("%s/cmd", t)
}

func (t DirPathFakeCommand) FilePathBehaviors() string {
	return fmt.Sprintf("%s/behaviors.json", t)
}

func (t DirPathFakeCommand) FilePathState() string {
	return fmt.Sprintf("%s/state.json", t)
}

func (t DirPathFakeCommand) FilePathProcessing() string {
	return fmt.Sprintf("%s/processing", t)
}

type DirPathFakeCommands []DirPathFakeCommand

type FakeCommand struct {
	filePathFakeCMD string
	dirPath         DirPathFakeCommand
	behaviors       Behaviors
}

func (t *FakeCommand) Init(force bool) error {
	if err := t.initBaseDir(force); err != nil {
		return err
	}

	if err := t.initCopyFakeCMD(); err != nil {
		return err
	}

	if err := t.initBehaviors(); err != nil {
		return err
	}

	return nil
}

func (t *FakeCommand) initBaseDir(force bool) error {
	if force {
		if err := os.RemoveAll(t.dirPath.String()); err != nil {
			return fmt.Errorf("failed to os.RemoveAll: %w", err)
		}
	}

	if err := os.MkdirAll(t.dirPath.String(), 0755); err != nil {
		return fmt.Errorf("failed to os.MkdirAll: %w", err)
	}

	return nil
}

func (t *FakeCommand) initCopyFakeCMD() error {
	src, err := os.Open(t.filePathFakeCMD)
	if err != nil {
		return fmt.Errorf("fakecmd is invalid: %s: %w", t.filePathFakeCMD, err)
	}
	defer src.Close()

	srcFI, err := src.Stat()
	if err != nil {
		return fmt.Errorf("failed to Stat: %w", err)
	}

	dstFilePathFakeCMD := t.dirPath.FilePathCommand()

	if _, err := os.Stat(dstFilePathFakeCMD); err == nil {
		return fmt.Errorf("cloned fakecmd already exists: %s", dstFilePathFakeCMD)
	}

	dst, err := os.Create(dstFilePathFakeCMD)
	if err != nil {
		return fmt.Errorf("failed to os.Create: %w", err)
	}
	defer dst.Close()

	if err := dst.Chmod(srcFI.Mode()); err != nil {
		return fmt.Errorf("failed to Chmod: %w", err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	return nil
}

func (t *FakeCommand) initBehaviors() error {
	b, err := json.MarshalIndent(t.behaviors, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to json.MarshalIndent: %w", err)
	}

	if err := os.WriteFile(t.dirPath.FilePathBehaviors(), b, 0644); err != nil {
		return fmt.Errorf("failed to os.WriteFile: %w", err)
	}

	return nil
}

func (t *FakeCommand) Cleanup() error {
	if err := os.RemoveAll(t.dirPath.String()); err != nil {
		return fmt.Errorf("failed to remove dir: %w", err)
	}
	return nil
}

func (t *FakeCommand) DirPath() DirPathFakeCommand {
	return t.dirPath
}

func NewFakeCommand(
	filePathFakeCMD string,
	dirPath DirPathFakeCommand,
	behaviors Behaviors,
) *FakeCommand {
	return &FakeCommand{
		filePathFakeCMD: filePathFakeCMD,
		dirPath:         dirPath,
		behaviors:       behaviors,
	}
}
