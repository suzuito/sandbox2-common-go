package e2ehelpers

import "net/http"

type RoundTripperForE2E struct {
	e2eTestID  string
	origin     http.RoundTripper
	fakeScheme string
	fakeHost   string
}

func (t *RoundTripperForE2E) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("E2E-TestId", t.e2eTestID)
	originalURL := req.URL
	originalURL.Scheme = t.fakeScheme
	originalURL.Host = t.fakeHost
	return t.origin.RoundTrip(req)
}

func NewRoundTripperForE2E(
	e2eTestID string,
	origin http.RoundTripper,
	fakeScheme string,
	fakeHost string,
) *RoundTripperForE2E {
	return &RoundTripperForE2E{
		e2eTestID:  e2eTestID,
		origin:     origin,
		fakeScheme: fakeScheme,
		fakeHost:   fakeHost,
	}
}
