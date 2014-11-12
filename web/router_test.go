package web_test

import (
	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web"
	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web/handlers"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Router", func() {
	var router web.Router

	BeforeEach(func() {
		router = web.NewRouter()
	})

	It("routes GET /info", func() {
		s := router.Routes().Get("GET /info").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetInfo{}))
	})

	It("routes POST /", func() {
		s := router.Routes().Get("POST /").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.ForwardEmail{}))
	})
})
