package e2ehelpers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smocker-dev/smocker/server/types"
	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
	"github.com/suzuito/sandbox2-common-go/libs/utils"
)

func Test_SmockerClient_PostMocks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc           string
		inputMocks     types.Mocks
		inputReset     bool
		wantErr        bool
		expectedErrMsg string
		server         *httptest.Server
	}{
		{
			desc:       "ok",
			inputMocks: types.Mocks{},
			inputReset: true,
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})),
		},
		{
			desc:       "ng - http error",
			inputMocks: types.Mocks{},
			inputReset: true,
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "TEST")
			})),
			wantErr:        true,
			expectedErrMsg: "http error: status=400 body=TEST",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			defer tC.server.Close()

			client := e2ehelpers.NewSmockerClient(
				utils.MustParseURL(tC.server.URL),
				http.DefaultClient,
			)

			err := client.PostMocks(tC.inputMocks, tC.inputReset)
			if tC.wantErr {
				require.EqualError(t, err, tC.expectedErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
