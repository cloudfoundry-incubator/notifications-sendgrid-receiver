package main

import "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/application"

func main() {
	app := application.NewApplication()
	defer app.Crash()

	app.Boot()
}
