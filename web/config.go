package web

import "log"

type Config struct {
	BasicAuthPassword string
	BasicAuthUsername string
	Logger            *log.Logger
	NotificationsHost string
	Port              string
	VerifySSL         bool
}
