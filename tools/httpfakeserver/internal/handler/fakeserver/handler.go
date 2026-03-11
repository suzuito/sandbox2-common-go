package fakeserver

import (
	"fmt"
	"net/http"

	"github.com/suzuito/sandbox2-common-go/tools/httpfakeserver/internal/domain/mock"
)

func HandleFunc(caseRepo *mock.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for m := range caseRepo.Mocks() {
			if !m.Match(r) {
				continue
			}

			m.WriteResponse(w)
			return
		}

		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "no matched to cases")
	}
}
