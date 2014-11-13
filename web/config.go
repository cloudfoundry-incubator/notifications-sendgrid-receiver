package web

import "log"

type Config struct {
	Port              string
	Logger            *log.Logger
	NotificationsHost string
	BasicAuthUsername string
	BasicAuthPassword string
}
