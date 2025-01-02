package e2ehelpers

import (
	"encoding/json"
	"os"
	"strings"
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
