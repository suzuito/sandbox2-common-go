package hfs

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/tools/httpfakeserver"
)

func mustJSONMarshal(t *testing.T, a any) []byte {
	b, err := json.Marshal(a)
	if err != nil {
		require.NoError(t, err)
	}
	return b
}

func mustJSONUnmarshalFromHTTPResponse[T any](t *testing.T, res *http.Response) *T {
	b, err := io.ReadAll(res.Body)
	if err != nil {
		require.NoError(t, err)
	}

	var a T
	if err := json.Unmarshal(b, &a); err != nil {
		require.NoError(t, err)
	}

	return &a
}

func TestNormalCases(t *testing.T) {
	cli := http.DefaultClient

	// setup: モックをクリアする
	req, err := http.NewRequest(http.MethodDelete, targetURL+"/admin/cases", nil)
	require.NoError(t, err)
	res, err := cli.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, res.StatusCode)
	require.NoError(t, res.Body.Close())

	// setup: 3つのモックを登録する(GET,POST,PUT)
	mocks := []httpfakeserver.Mock{
		{
			Request: httpfakeserver.Request{
				Method: "GET",
				Path:   "/foo",
				Header: http.Header{"X-Custom": []string{"val1"}},
				Query:  url.Values{"key": []string{"val1"}},
			},
			Response: httpfakeserver.Response{
				Status: http.StatusOK,
				Body:   `{"message":"get-foo"}`,
				Header: http.Header{"Content-Type": []string{"application/json"}},
			},
		},
		{
			Request: httpfakeserver.Request{
				Method: "POST",
				Path:   "/bar",
				Header: http.Header{"X-Custom": []string{"val2"}},
				Query:  url.Values{"key": []string{"val2"}},
			},
			Response: httpfakeserver.Response{
				Status: http.StatusCreated,
				Body:   `{"message":"post-bar"}`,
				Header: http.Header{"Content-Type": []string{"application/json"}},
			},
		},
		{
			Request: httpfakeserver.Request{
				Method: "PUT",
				Path:   "/baz",
				Header: http.Header{"X-Custom": []string{"val3"}},
				Query:  url.Values{"key": []string{"val3"}},
			},
			Response: httpfakeserver.Response{
				Status: http.StatusAccepted,
				Body:   `{"message":"put-baz"}`,
				Header: http.Header{"Content-Type": []string{"application/json"}},
			},
		},
	}

	for _, m := range mocks {
		res, err = cli.Post(
			targetURL+"/admin/cases", "application/json",
			bytes.NewBuffer(mustJSONMarshal(t, m)),
		)
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, res.StatusCode)
		require.NoError(t, res.Body.Close())
	}

	// case: 1つ目のモックに一致するリクエストを投げる -> モックで設定したレスポンスが返る
	t.Run("matched request returns mock response", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, targetURL+"/foo?key=val1", nil)
		require.NoError(t, err)
		req.Header.Set("X-Custom", "val1")

		res, err := cli.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, `{"message":"get-foo"}`, string(body))
	})

	// case: Methodだけ異なる -> 501
	t.Run("different method returns 501", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, targetURL+"/foo?key=val1", nil)
		require.NoError(t, err)
		req.Header.Set("X-Custom", "val1")

		res, err := cli.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotImplemented, res.StatusCode)
	})

	// case: Headerだけ異なる -> 501
	t.Run("different header returns 501", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, targetURL+"/foo?key=val1", nil)
		require.NoError(t, err)
		req.Header.Set("X-Custom", "wrong")

		res, err := cli.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotImplemented, res.StatusCode)
	})

	// case: Queryだけ異なる -> 501
	t.Run("different query returns 501", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, targetURL+"/foo?key=wrong", nil)
		require.NoError(t, err)
		req.Header.Set("X-Custom", "val1")

		res, err := cli.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotImplemented, res.StatusCode)
	})
}

func TestEdgeCases(t *testing.T) {
	cli := http.DefaultClient

	// setup: モックをクリアする
	req, err := http.NewRequest(http.MethodDelete, targetURL+"/admin/cases", nil)
	require.NoError(t, err)
	res, err := cli.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, res.StatusCode)
	require.NoError(t, res.Body.Close())

	t.Run("unregistered request returns 501", func(t *testing.T) {
		res, err := cli.Get(targetURL + "/abcde")
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotImplemented, res.StatusCode)
	})

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		res, err := cli.Post(
			targetURL+"/admin/cases", "application/json",
			bytes.NewBufferString("this is not json"),
		)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("empty mock can be registered and retrieved", func(t *testing.T) {
		res, err := cli.Post(
			targetURL+"/admin/cases", "application/json",
			bytes.NewBuffer(mustJSONMarshal(
				t,
				httpfakeserver.Mock{},
			)),
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		require.NoError(t, res.Body.Close())

		res, err = cli.Get(targetURL + "/admin/cases")
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		actual := mustJSONUnmarshalFromHTTPResponse[httpfakeserver.Mocks](t, res)
		assert.Equal(
			t,
			&httpfakeserver.Mocks{httpfakeserver.Mock{}},
			actual,
		)
	})
}
