package uaa

import (
    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/config"
    uaa_lib "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type UAAClient struct {
    auth uaa_lib.UAA
}

type UAAClientInterface interface {
    AccessToken() string
}

func NewUAAClient(env config.Environment) UAAClient {
    auth := uaa_lib.NewUAA("", env.UAAHost, env.UAAClientID, env.UAAClientSecret, "")
    auth.VerifySSL = env.VerifySSL

    return UAAClient{
        auth: auth,
    }
}

func (client UAAClient) retrieveUAAClientToken() uaa_lib.Token {
    token, err := client.auth.GetClientToken()
    if err != nil {
        panic(err)
    }
    return token
}

func (client UAAClient) AccessToken() string {
    token := client.retrieveUAAClientToken()
    return token.Access
}
