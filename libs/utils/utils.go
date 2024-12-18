package utils

import "net/url"

func Ptr[T any](i T) *T {
	return &i
}

func MustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}
