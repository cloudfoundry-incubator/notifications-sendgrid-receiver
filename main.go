package main

import "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web"

func main() {
    server := web.NewServer()
    server.Run()
}
