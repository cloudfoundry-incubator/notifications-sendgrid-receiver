package web_test

import (
	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/fakes"
	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web"
	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web/handlers"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	var router *mux.Router

	BeforeEach(func() {
		uaaClient := fakes.NewUAAClient()
		router = web.NewServer(web.Config{}).Router(uaaClient)
	})

	It("routes GET /info", func() {
		s := router.Get("GET /info").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetInfo{}))
	})

	It("routes POST /", func() {
		s := router.Get("POST /").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.ForwardEmail{}))
	})
})
