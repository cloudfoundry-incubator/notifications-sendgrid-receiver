package requests

import (
    "encoding/base64"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/config"
)

type BasicAuthenticatorInterface interface {
    Verify(header http.Header) bool
}

type BasicAuthenticator struct{}

func NewBasicAuthenticator() BasicAuthenticator {
    return BasicAuthenticator{}
}

func (authenticator BasicAuthenticator) Verify(header http.Header) bool {
    authorization := header.Get("Authorization")
    basicAuthToken := strings.TrimPrefix(authorization, "Basic ")
    credentials, err := base64.StdEncoding.DecodeString(basicAuthToken)
    if err != nil {
        return false
    }

    env := config.NewEnvironment()
    usernamePassword := strings.Split(string(credentials), ":")

    if len(usernamePassword) == 2 {
        if usernamePassword[0] == env.BasicAuthUserName && usernamePassword[1] == env.BasicAuthPassword {
            return true
        }
    }

    return false
}
