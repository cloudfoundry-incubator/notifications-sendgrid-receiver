package application

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web"
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

	logger.Println("Trying to boot")

	logger.Println("Booting with configuration:")
	viron.Print(env, logger)

	server := web.NewServer(web.Config{
		Port:              env.Port,
		Logger:            logger,
		NotificationsHost: env.NotificationsHost,
		BasicAuthUsername: env.BasicAuthUsername,
		BasicAuthPassword: env.BasicAuthPassword,
	})
	server.Run(mother.UAAClient())
}

// This is a hack to get the logs to output to the loggregator before the process exits
func (app Application) Crash() {
	err := recover()
	if err != nil {
		time.Sleep(5 * time.Second)
		panic(err)
	}
}
