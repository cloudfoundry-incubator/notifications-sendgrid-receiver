package services

import (
	"encoding/base64"
	"net/http"
	"strings"
)

type BasicAuthenticator struct {
	username string
	password string
}

func NewBasicAuthenticator(username, password string) BasicAuthenticator {
	return BasicAuthenticator{
		username: username,
		password: password,
	}
}

func (authenticator BasicAuthenticator) Verify(header http.Header) bool {
	credentials, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(header.Get("Authorization"), "Basic "))
	if err != nil {
		return false
	}

	usernamePassword := strings.Split(string(credentials), ":")
	if len(usernamePassword) == 2 {
		if usernamePassword[0] == authenticator.username && usernamePassword[1] == authenticator.password {
			return true
		}
	}

	return false
}
