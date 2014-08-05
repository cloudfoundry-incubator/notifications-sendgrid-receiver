package web

import (
    "strings"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/requests"
    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web/handlers"
    "github.com/gorilla/mux"
    "github.com/ryanmoran/stack"
)

type Router struct {
    stacks map[string]stack.Stack
}

func NewRouter() Router {
    return Router{
        stacks: map[string]stack.Stack{
            "GET /info": stack.NewStack(handlers.NewGetInfo()),
            "POST /":    stack.NewStack(handlers.NewForwardEmail(requests.NewRequestBuilder(), requests.NewRequestSender())),
        },
    }
}

func (router Router) Routes() *mux.Router {
    r := mux.NewRouter()
    for methodPath, stack := range router.stacks {
        var name = methodPath
        parts := strings.SplitN(methodPath, " ", 2)
        r.Handle(parts[1], stack).Methods(parts[0]).Name(name)
    }
    return r
}
