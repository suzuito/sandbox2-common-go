package admin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/suzuito/sandbox2-common-go/tools/httpfakeserver/internal/domain/mock"
)

func PostAdminCase(caseRepo *mock.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "failed to read body") //nolint:errcheck
			return
		}

		c := mock.Mock{}
		if err := json.Unmarshal(body, &c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "failed to parse body") //nolint:errcheck
			return
		}

		caseRepo.SetMock(c)

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteAdminCase(caseRepo *mock.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		caseRepo.Clear()

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetAdminCase(caseRepo *mock.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ret := []mock.Mock{}
		for m := range caseRepo.Mocks() {
			ret = append(ret, m)
		}

		body, err := json.Marshal(ret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "failed to marshal mocks") //nolint:errcheck
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body) // nolint:errcheck
	}
}
