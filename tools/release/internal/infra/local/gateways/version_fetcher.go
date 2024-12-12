package gateways

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/gateways"
)

type VersionFetcher struct {
	filePathGit string
}

var _ gateways.VersionFetcher = &VersionFetcher{}

func (t *VersionFetcher) GetLatestVersion(ctx context.Context, prefix string) (*semver.Version, error) {
	cmdResultString := bytes.NewBufferString("")
	cmd := exec.CommandContext(
		ctx,
		t.filePathGit,
		"tag",
	)
	cmd.Stdout = cmdResultString

	if err := cmd.Run(); err != nil {
		var exiterr *exec.ExitError
		if errors.As(err, &exiterr) {
			return nil, fmt.Errorf("'git tag' command is failed with code %d: %w", exiterr.ExitCode(), err)
		}
		return nil, fmt.Errorf("'git tag' commmand is failed: %w", err)
	}

	cmdResultLines := strings.Split(cmdResultString.String(), "\n")

	versions := make([]*semver.Version, 0)
	for _, line := range cmdResultLines {
		versionString, _ := strings.CutPrefix(line, prefix)

		version, err := semver.StrictNewVersion(versionString)
		if err != nil {
			continue
		}

		versions = append(versions, version)
	}

	if len(versions) <= 0 {
		return nil, fmt.Errorf("version string does not exist in git tags")
	}

	sort.Sort(semver.Collection(versions))

	return versions[len(versions)-1], nil
}

func NewVersionFetcher(
	filePathGit string,
) *VersionFetcher {
	return &VersionFetcher{
		filePathGit: filePathGit,
	}
}
