package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Forward", func() {
    var handler handlers.ForwardEmail
    var fakeUAAServer *httptest.Server
    var fakeRequestBuilder FakeRequestBuilder
    var fakeRequestSender FakeRequestSender

    BeforeEach(func() {
        fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            response := map[string]string{
                "access_token": "fakeAuthToken",
            }
            responseBody, err := json.Marshal(response)
            if err != nil {
                panic(err)
            }
            w.Write(responseBody)
        }))

        err := os.Setenv("UAA_HOST", fakeUAAServer.URL)
        if err != nil {
            panic(err)
        }

        err = os.Setenv("UAA_CLIENT_ID", "fake_client_id")
        if err != nil {
            panic(err)
        }

        err = os.Setenv("UAA_CLIENT_SECRET", "fake_client_secret")
        if err != nil {
            panic(err)
        }

        handler = handlers.NewForwardEmail(&fakeRequestBuilder, &fakeRequestSender)
    })

    AfterEach(func() {
        defer fakeUAAServer.Close()
    })

    Describe("ServeHTTP", func() {
        It("sends a request built by the request builder  to the notifications service", func() {
            writer := httptest.NewRecorder()
            body, err := json.Marshal(map[string]string{
                "stuff": "lol",
            })

            request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
            if err != nil {
                panic(err)
            }

            handler.ServeHTTP(writer, request)

            Expect(fakeRequestSender.Request).To(Equal(fakeRequestBuilder.Request))
        })

        Context("when notifications responds with a success", func() {
            // TODO test that we hit UAA for access token

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
        })
    })
})
