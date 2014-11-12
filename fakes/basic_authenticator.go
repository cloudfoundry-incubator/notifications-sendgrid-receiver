package fakes

import "net/http"

type BasicAuthenticator struct {
	InvalidAuth bool
}

func NewBasicAuthenticator() *BasicAuthenticator {
	return &BasicAuthenticator{}
}

func (fake BasicAuthenticator) Verify(header http.Header) bool {
	if fake.InvalidAuth {
		return false
	}
	return true
}
