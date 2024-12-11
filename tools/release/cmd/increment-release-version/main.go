package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
	var incrementTypeString string
	var prefix string
	flag.StringVar(
		&prefix,
		"prefix",
		"",
		"バージョン文字列のプレフィクス",
	)
	flag.StringVar(
		&incrementTypeString,
		"increment",
		"patch",
		"どのバージョンをインクリメントするかを指定する(major,minor or patch, default: patch)",
	)
	flag.Usage = usage

	flag.Parse()

	incrementType := IncrementType(incrementTypeString)
	if err := incrementType.Validate(); err != nil {
		log.Fatalf("invalid increment type '%s'", incrementType)
	}

	//	version, err := semver.StrictNewVersion(versionString)
	//	if err != nil {
	//		log.Fatalf("cannot parse '%s' as a semantic version", versionString)
	//	}
}
