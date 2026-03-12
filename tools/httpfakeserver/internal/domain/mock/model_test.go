package mock

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestEqual(t *testing.T) {
	assertEqualRequestID(
		t,
		Request{},
		Request{},
	)
	assertEqualRequestID(
		t,
		Request{
			Method: http.MethodGet,
		},
		Request{
			Method: http.MethodGet,
		},
	)
}

func assertEqualRequestID(t *testing.T, l, r Request) {
	lid := l.ID()
	rid := r.ID()
	assert.Equal(t, lid, rid)
}
