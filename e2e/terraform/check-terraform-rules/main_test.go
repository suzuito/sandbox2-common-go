package main

import (
	"fmt"
	"os"
	"testing"
)

func TestA(t *testing.T) {
	filePathBin := os.Getenv("FILE_PATH_BIN")

	fmt.Println(filePathBin)
}
