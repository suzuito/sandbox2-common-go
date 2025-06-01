package errordefcli_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	errordefcli "github.com/suzuito/sandbox2-common-go/libs/errordefs/cli"
)

func TestErrordefcli(t *testing.T) {
	type arg struct {
		newFunc     func() error
		defaultCode uint8
	}

	testCases := []struct {
		desc                string
		arg                 arg
		expectedCodeCode    uint8
		expectedCodeMessage string
	}{
		{
			desc: "ok - NewCLIErrorf",
			arg: arg{
				newFunc: func() error {
					return errordefcli.NewCLIErrorf(1, "dummy1")
				},
			},
			expectedCodeCode:    1,
			expectedCodeMessage: "cli error: dummy1",
		},
		{
			desc: "ok - NewCLIErrorf - can use %%w",
			arg: arg{
				newFunc: func() error {
					return errordefcli.NewCLIErrorf(1, "dummy1: %w", os.ErrExist)
				},
			},
			expectedCodeCode:    1,
			expectedCodeMessage: "cli error: dummy1: file already exists",
		},
		{
			arg: arg{
				newFunc: func() error {
					return errors.New("not cli error")
				},
				defaultCode: 2,
			},
			expectedCodeCode:    2,
			expectedCodeMessage: "not cli error",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.arg.newFunc()
			actualCodeCode, actualCodeMessage := errordefcli.Code(err, tC.arg.defaultCode)
			assert.Equal(t, tC.expectedCodeCode, actualCodeCode)
			assert.Equal(t, tC.expectedCodeMessage, actualCodeMessage)
		})
	}
}

func TestErrordefcliUnwrap(t *testing.T) {
	err := errordefcli.NewCLIErrorf(1, "dummy: %w", os.ErrExist)
	assert.True(t, errors.Is(err, os.ErrExist))
}
