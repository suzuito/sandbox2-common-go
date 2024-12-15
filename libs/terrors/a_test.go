package terrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	var exact error
	exact = Wrap(nil)
	assert.Nil(t, exact)
	exact = Wrap(fmt.Errorf("dummy error"))
	assert.Regexp(t, "^dummy error", exact.Error())
	assert.Regexp(t, ".+/a_test.go", exact.(*traceableErrorImpl).TraceInfo.Filename)
	assert.Equal(t, 15, exact.(*traceableErrorImpl).TraceInfo.Line)
}

func TestStackTrace(t *testing.T) {
	err1 := fmt.Errorf("err1")
	err2 := Wrap(err1)
	err3 := fmt.Errorf("%w", err2)
	err4 := Wrap(err3)
	traceInfos := err2.(*traceableErrorImpl).StackTrace()
	assert.Equal(t, 1, len(traceInfos))
	assert.Regexp(t, ".+/a_test.go", traceInfos[0].Filename)
	assert.Equal(t, 23, traceInfos[0].Line)
	traceInfos = err4.(*traceableErrorImpl).StackTrace()
	assert.Equal(t, 2, len(traceInfos))
	assert.Regexp(t, ".+/a_test.go", traceInfos[0].Filename)
	assert.Equal(t, 25, traceInfos[0].Line)
	assert.Regexp(t, ".+/a_test.go", traceInfos[1].Filename)
	assert.Equal(t, 23, traceInfos[1].Line)

	err5 := Wrapf("this is a test error: %w", err1)
	assert.Equal(t, "this is a test error: err1", err5.Error())
	assert.Equal(t, "err1", errors.Unwrap(errors.Unwrap(err5)).Error())
}
