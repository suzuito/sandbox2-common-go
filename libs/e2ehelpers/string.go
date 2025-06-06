package e2ehelpers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// NewLines returns joined string by "\n"
func NewLines(a ...string) string {
	return strings.Join(a, "\n") + "\n"
}

func MinifyJSONString(a string) string {
	v := map[string]any{}
	if err := json.Unmarshal([]byte(a), &v); err != nil {
		panic(err)
	}

	r, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(r)
}

func MustMarshalJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return b
}

func MustWriteJSONFile(filePath string, v any) {
	MustWriteFile(filePath, MustMarshalJSON(v))
}

func MustWriteFile(filePath string, b []byte) {
	if err := os.WriteFile(filePath, b, 0755); err != nil {
		panic(err)
	}
}

func MustMkdir(dirPath string) {
	if err := os.Mkdir(dirPath, 0755); err != nil {
		panic(err)
	}
}

func MustWriteFileAtRandomPath(baseDir string, b []byte) string {
	fp := filepath.Join(baseDir, uuid.NewString())
	MustWriteFile(fp, b)
	return fp
}
