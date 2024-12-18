package e2ehelpers

import "strings"

// NewLines returns joined string by "\n"
func NewLines(a ...string) string {
	return strings.Join(a, "\n") + "\n"
}
