package web

import (
    "log"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/config"
)

type Server struct {
}

func NewServer() Server {
    return Server{}
}

func (s Server) Run() {
    env := config.NewEnvironment()

    router := NewRouter()
    log.Printf("Listening on localhost:%s\n", env.Port)

    http.ListenAndServe(":"+env.Port, router.Routes())
}
