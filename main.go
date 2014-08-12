package main

import (
    "time"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/config"
    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/log"
    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web"
)

func main() {
    defer crash()
    log.Println("Trying to boot")

    env := config.NewEnvironment()
    configure(env)
    server := web.NewServer()
    server.Run()
}

func configure(env config.Environment) {
    log.Println("Booting with configuration:")
    log.Printf("\tUAAHost                 -> %+v", env.UAAHost)
    log.Printf("\tCCHost                  -> %+v", env.CCHost)
    log.Printf("\tVerifySSL               -> %+v", env.VerifySSL)
    log.Printf("\tNotificationsHost       -> %+v", env.NotificationsHost)
    log.Printf("\tPort                    -> %+v", env.Port)
}

// This is a hack to get the logs to output to the loggregator before the process exits
func crash() {
    err := recover()
    if err != nil {
        time.Sleep(5 * time.Second)
        panic(err)
    }
}
