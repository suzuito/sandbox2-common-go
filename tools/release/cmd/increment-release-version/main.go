package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	errordefcli "github.com/suzuito/sandbox2-common-go/libs/errordefs/cli"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/inject"
)

var usageString = `git tagsコマンドを実行したリリースバージョン文字列を検証します。下記を検証しバージョン文字列が誤っている場合、エラーを出力し異常終了します。

* リリースバージョン文字列がセマンティックバージョンに準拠しているか。https://semver.org/
* リリースバージョン文字列が既存の最新バージョンよりも新しいか。

`

func usage() {
	fmt.Fprintln(os.Stderr, usageString)
	flag.PrintDefaults()
}

func main() {
	ctx := context.Background()

	var incrementType string
	var prefix string
	var filePathGit string
	var githubOwner string
	var githubRepo string
	var githubAppToken string
	var branch string
	flag.StringVar(
		&prefix,
		"prefix",
		"",
		"prefix of version string",
	)
	flag.StringVar(
		&incrementType,
		"increment",
		"patch",
		"which version is incremented(major,minor or patch, default: patch)",
	)
	flag.StringVar(
		&filePathGit,
		"git",
		"",
		"file path of git binary",
	)
	flag.StringVar(
		&githubOwner,
		"owner",
		"",
		"github owner",
	)
	flag.StringVar(
		&githubRepo,
		"repo",
		"",
		"github repo",
	)
	flag.StringVar(
		&branch,
		"branch",
		"",
		"github branch",
	)
	flag.StringVar(
		&githubAppToken,
		"token",
		"",
		"github app token",
	)
	flag.Usage = usage

	flag.Parse()

	if filePathGit == "" {
		fmt.Fprintf(os.Stderr, "-git is required\n")
		os.Exit(1)
	}
	if githubOwner == "" {
		fmt.Fprintf(os.Stderr, "-owner is required\n")
		os.Exit(1)
	}
	if githubRepo == "" {
		fmt.Fprintf(os.Stderr, "-repo is required\n")
		os.Exit(1)
	}
	if branch == "" {
		fmt.Fprintf(os.Stderr, "-branch is required\n")
		os.Exit(1)
	}
	if githubAppToken == "" {
		fmt.Fprintf(os.Stderr, "-token is required\n")
		os.Exit(1)
	}

	uc, err := inject.NewUsecase(
		filePathGit,
		githubAppToken,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to inject.NewUsecase: %+v\n", err)
		os.Exit(1)
	}

	if err := uc.IncrementVersion(
		ctx,
		githubOwner,
		githubRepo,
		branch,
		prefix,
		incrementType,
	); err != nil {
		if clierr, ok := errordefcli.AsCLIError(err); ok {
			fmt.Fprintln(os.Stderr, clierr.Error())
			os.Exit(clierr.ExitCode())
		}

		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(125)
	}
}
