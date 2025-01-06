package e2ehelpers_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
)

func TestTestID(t *testing.T) {
	tid := e2ehelpers.NewTestID()
	require.Equal(t, tid.UUID().String(), tid.String())
}
