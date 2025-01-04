package domains

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type Faker struct {
	filePathFakeCMD string

	dirPaths DirPathFakeCommands
}

func (t *Faker) Add(
	dirPathCommand DirPathFakeCommand,
	behaviors Behaviors,
) *FakeCommand {
	t.dirPaths = append(t.dirPaths, dirPathCommand)

	fcmd := FakeCommand{
		filePathFakeCMD: t.filePathFakeCMD,
		dirPath:         dirPathCommand,
		behaviors:       behaviors,
	}

	return &fcmd
}

func (t *Faker) Cleanup() error {
	errs := []error{}
	for _, d := range t.dirPaths {
		errs = append(errs, os.RemoveAll(d.String()))
	}
	return errors.Join(errs...)
}

func (t *Faker) AddInTest(tt *testing.T, behaviors Behaviors) *FakeCommand {
	fcmd := t.Add(
		DirPathFakeCommand(fmt.Sprintf("/tmp/%s", uuid.NewString())),
		behaviors,
	)
	require.NoError(tt, fcmd.Init(false))
	return fcmd
}

func New(
	filePathFakeCMD string,
) *Faker {
	return &Faker{
		filePathFakeCMD: filePathFakeCMD,
		dirPaths:        DirPathFakeCommands{},
	}
}

const envName = "FILE_PATH_FAKECMD"

func MustByEnv() *Faker {
	envVal := os.Getenv(envName)
	if envVal == "" {
		panic(fmt.Errorf("environment variable '%s' is empty", envName))
	}

	return New(envVal)
}
