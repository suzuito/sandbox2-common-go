package mock

import (
	"fmt"
	"iter"
	"maps"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Method string      `json:"method"`
	Path   string      `json:"path"`
	Header http.Header `json:"header"`
	Query  url.Values  `json:"query"`
}

func (r *Request) ID() string {
	header := strings.Builder{}
	r.Header.Write(&header) // nolint:errcheck
	src := strings.ToLower(r.Method) + r.Path + header.String() + r.Query.Encode()
	return src
}

func (r *Request) Equal(rr *Request) bool {
	return r.ID() == rr.ID()
}

func NewRequestFromHTTPRequest(
	r *http.Request,
	includeHeaderKeys iter.Seq[string],
) *Request {
	rr := Request{}
	rr.Method = r.Method
	rr.Path = r.URL.Path
	rr.Header = http.Header{}
	rr.Query = url.Values{}
	for k := range includeHeaderKeys {
		vs := r.Header[k]
		for _, v := range vs {
			rr.Header.Set(k, v)
		}
	}
	for k, vs := range r.URL.Query() {
		for _, v := range vs {
			rr.Query.Set(k, v)
		}
	}
	return &rr
}

type Response struct {
	Header http.Header `json:"header"`
	Body   string      `json:"body"`
	Status int         `json:"status"`
}

type Mock struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

func (c Mock) ID() string {
	return c.Request.ID()
}

func (c Mock) Match(r *http.Request) bool {
	headerKeys := maps.Keys(c.Request.Header)
	rr := NewRequestFromHTTPRequest(r, headerKeys)
	return c.Request.Equal(rr)
}

func (c Mock) WriteResponse(w http.ResponseWriter) {
	for k, vs := range c.Response.Header {
		for _, v := range vs {
			w.Header().Set(k, v)
		}
	}
	w.WriteHeader(c.Response.Status)
	fmt.Fprint(w, c.Response.Body) //nolint:errcheck
}

type Mocks []Mock
