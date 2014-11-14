package web

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web/handlers"
	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web/services"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type Server struct {
	config Config
}

func NewServer(config Config) Server {
	return Server{
		config: config,
	}
}

func (server Server) Run(uaaClient handlers.UAAClientInterface) {
	server.config.Logger.Printf("Listening on localhost:%s\n", server.config.Port)

	http.ListenAndServe(":"+server.config.Port, server.Router(uaaClient))
}

func (server Server) Router(uaaClient handlers.UAAClientInterface) *mux.Router {
	requestBuilder := services.NewRequestBuilder(server.config.NotificationsHost)
	requestSender := services.NewRequestSender(server.config.Logger, server.config.VerifySSL)
	requestParser := services.NewRequestBodyParser()
	basicAuthenticator := services.NewBasicAuthenticator(server.config.BasicAuthUsername, server.config.BasicAuthPassword)

	stacks := map[string]stack.Stack{
		"GET /info": stack.NewStack(handlers.NewGetInfo()),
		"POST /":    stack.NewStack(handlers.NewForwardEmail(requestBuilder, requestSender, uaaClient, requestParser, basicAuthenticator, server.config.Logger)),
	}

	r := mux.NewRouter()
	for methodPath, stack := range stacks {
		var name = methodPath
		parts := strings.SplitN(methodPath, " ", 2)
		r.Handle(parts[1], stack).Methods(parts[0]).Name(name)
	}

	return r
}
