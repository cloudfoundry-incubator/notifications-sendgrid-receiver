package uaa

import (
	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/config"
	uaa_lib "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type UAAClient struct {
	auth uaa_lib.UAA
}

type UAAClientInterface interface {
	AccessToken() (string, error)
}

func NewUAAClient(env config.Environment) UAAClient {
	auth := uaa_lib.NewUAA("", env.UAAHost, env.UAAClientID, env.UAAClientSecret, "")
	auth.VerifySSL = env.VerifySSL

	return UAAClient{
		auth: auth,
	}
}

func (client UAAClient) retrieveUAAClientToken() (uaa_lib.Token, error) {
	return client.auth.GetClientToken()
}

func (client UAAClient) AccessToken() (string, error) {
	token, err := client.retrieveUAAClientToken()
	return token.Access, err
}
