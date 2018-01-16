package application

import "github.com/ryanmoran/viron"

var UAAPublicKey string

type Environment struct {
	BasicAuthPassword string `env:"BASIC_AUTH_PASSWORD"`
	BasicAuthUsername string `env:"BASIC_AUTH_USER_NAME"`
	CCHost            string `env:"CC_HOST"            env-required:"true"`
	LogFile           string `env:"LOG_FILE"`
	MaxRequests       int    `env:"MAX_REQUESTS"       env-default:"2"`
	NotificationsHost string `env:"NOTIFICATIONS_HOST" env-required:"true"`
	Port              string `env:"PORT"               env-required:"true"`
	UAAClientID       string `env:"UAA_CLIENT_ID"      env-required:"true"`
	UAAClientSecret   string `env:"UAA_CLIENT_SECRET"  env-required:"true"`
	UAAHost           string `env:"UAA_HOST"           env-required:"true"`
	VerifySSL         bool   `env:"VERIFY_SSL"         env-default:"true"`
}

func NewEnvironment() Environment {
	env := Environment{}
	err := viron.Parse(&env)
	if err != nil {
		panic(err)
	}

	return env
}
