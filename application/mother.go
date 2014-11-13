package application

import (
	"log"
	"os"

	"github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type Mother struct {
	logger *log.Logger
}

func NewMother() *Mother {
	return &Mother{}
}

func (mother *Mother) Logger() *log.Logger {
	if mother.logger == nil {
		mother.logger = log.New(os.Stdout, "", 0)
	}

	return mother.logger
}

func (mother *Mother) Environment() Environment {
	return NewEnvironment()
}

func (mother *Mother) UAAClient() uaa.UAA {
	env := mother.Environment()
	client := uaa.NewUAA("", env.UAAHost, env.UAAClientID, env.UAAClientSecret, "")
	client.VerifySSL = env.VerifySSL

	return client
}
