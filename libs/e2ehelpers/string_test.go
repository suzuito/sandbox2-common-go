package e2ehelpers_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
)

func TestNewLines(t *testing.T) {
	require.Equal(t, "hoge\nfuga\n", e2ehelpers.NewLines("hoge", "fuga"))
}
