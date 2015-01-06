package application

import (
	"time"

	"github.com/pivotal-cf/uaa-sso-golang/uaa"
	"github.com/ryanmoran/viron"
)

type Application struct{}

func NewApplication() Application {
	return Application{}
}

func (app Application) Boot() {
	mother := NewMother()
	logger := mother.Logger()
	env := mother.Environment()
	server := mother.Server()
	uaaClient := mother.UAAClient()
	app.register(uaaClient, mother.Registrar())

	logger.Println("Booting with configuration:")
	viron.Print(env, logger)

	server.Run(uaaClient)
}

func (app Application) register(uaaClient uaa.UAA, registrar Registrar) {
	token, err := uaaClient.GetClientToken()
	if err != nil {
		panic(err)
	}

	err = registrar.Register(token.Access)
	if err != nil {
		panic(err)
	}
}

// This is a hack to get the logs to output to the loggregator before the process exits
func (app Application) Crash() {
	err := recover()
	if err != nil {
		time.Sleep(5 * time.Second)
		panic(err)
	}
}
