package handlers

import (
	"net/http"

	"github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type UAAClientInterface interface {
	uaa.GetClientTokenInterface
}

type RequestBuilderInterface interface {
	Build(map[string]string, string) (*http.Request, error)
}

type RequestSenderInterface interface {
	Send(*http.Request) error
}

type RequestBodyParserInterface interface {
	Parse(*http.Request) (map[string]string, error)
}

type BasicAuthenticatorInterface interface {
	Verify(header http.Header) bool
}

type ThrottlerInterface interface {
	Throttle() bool
	Finish()
}
