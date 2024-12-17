package e2ehelpers

import "strings"

func NewLines(a ...string) string {
	return strings.Join(a, "\n") + "\n"
}
