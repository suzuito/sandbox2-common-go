package e2ehelpers

import (
	"encoding/json"
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
