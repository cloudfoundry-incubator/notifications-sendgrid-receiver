package application

import (
	"time"

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

	logger.Println("Booting with configuration:")
	viron.Print(env, logger)

	server.Run(uaaClient)
}

// This is a hack to get the logs to output to the loggregator before the process exits
func (app Application) Crash() {
	err := recover()
	if err != nil {
		time.Sleep(5 * time.Second)
		panic(err)
	}
}
