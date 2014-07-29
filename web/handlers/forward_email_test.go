package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Forward", func() {
    Describe("ServeHTTP", func() {
        var handler handlers.ForwardEmail
        fakeSpaceMailerApi := &handlers.FakeSpaceMailerApi{}

        BeforeEach(func() {
            handler = handlers.NewForwardEmail(fakeSpaceMailerApi)
        })

        It("returns a 200 response code and an empty JSON body", func() {
            writer := httptest.NewRecorder()
            request, err := http.NewRequest("POST", "/", nil)
            if err != nil {
                panic(err)
            }

            handler.ServeHTTP(writer, request)

            Expect(writer.Code).To(Equal(http.StatusOK))
            Expect(writer.Body.String()).To(Equal("{}"))
        })

        It("send a notification to the notifications service with space guid", func() {
            writer := httptest.NewRecorder()
            body, err := json.Marshal(map[string]string{
                "headers": "horseman",
                "text":    "Where's my head?",
                "html":    "<b>Where's my head!</b>",
                "from":    "horseman@example.com",
                "to":      "space-guid-foo123-bar456@example.com",
                "cc":      "johnnydepp@example.com",
                "subject": "Banana Damage",
            })

            request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
            if err != nil {
                panic(err)
            }

            handler.ServeHTTP(writer, request)

            Expect(fakeSpaceMailerApi.SpaceGuid).To(Equal("foo123-bar456"))
            Expect(fakeSpaceMailerApi.Params).To(Equal("blank"))
        })
    })
})
