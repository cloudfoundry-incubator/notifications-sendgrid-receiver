package handlers_test

import (
    "errors"
    "net/http"
    "net/http/httptest"
    "os"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web/handlers"

    "github.com/gorilla/mux"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Notifications Receiver", func() {
    Describe("#PostToSpace", func() {
        var receiver handlers.NotificationsReceiver
        BeforeEach(func() {
            receiver = handlers.NotificationsReceiver{}

        })

        It("sets the required Authorization and Content-Type HTTP Headers", func() {
            var authorizationHeader string
            var contentTypeHeader string

            fakeRouter := mux.NewRouter()
            fakeRouter.HandleFunc("/spaces/{guid}", func(w http.ResponseWriter, req *http.Request) {
                authorizationHeader = req.Header.Get("Authorization")
                contentTypeHeader = req.Header.Get("Content-Type")
            }).Methods("POST")

            fakeNotificationsServer := httptest.NewServer(fakeRouter)
            defer fakeNotificationsServer.Close()
            os.Setenv("NOTIFICATIONS_HOST", fakeNotificationsServer.URL)

            params := map[string]string{"to": "space-guid-mammoth1-banana2-damage3@example.com"}

            receiver.PostToSpace("fakeTokenThatNeedsToBeCreatedOrSomething", params)

            Expect(contentTypeHeader).To(Equal("application/x-www-form-urlencoded"))
            Expect(authorizationHeader).To(Equal("Bearer fakeTokenThatNeedsToBeCreatedOrSomething"))
        })

        It("makes a post request to a notification sender service /spaces/{guid}", func() {
            var guid, kind, text string

            fakeRouter := mux.NewRouter()
            fakeRouter.HandleFunc("/spaces/{guid}", func(w http.ResponseWriter, req *http.Request) {
                err := req.ParseForm()
                if err != nil {
                    panic(err)
                }

                guid = mux.Vars(req)["guid"]
                kind = req.Form.Get("kind")
                text = req.Form.Get("text")
            }).Methods("POST")

            fakeNotificationsServer := httptest.NewServer(fakeRouter)
            defer fakeNotificationsServer.Close()
            os.Setenv("NOTIFICATIONS_HOST", fakeNotificationsServer.URL)

            params := map[string]string{
                "to":   "space-guid-mammoth1-banana2-damage3@example.com",
                "kind": "spacemail",
                "text": "Contents of the email message",
            }

            receiver.PostToSpace("fakeTokenThatNeedsToBeCreatedOrSomething", params)

            Expect(guid).To(Equal("mammoth1-banana2-damage3"))
            Expect(kind).To(Equal("spacemail"))
            Expect(text).To(Equal("Contents of the email message"))
        })

        It("returns an error if we are unable to parse the space guid", func() {
            params := map[string]string{
                "to": "",
            }
            err := receiver.PostToSpace("fakeTokenThatNeedsToBeCreatedOrSomething", params)
            Expect(err).To(Equal(errors.New("Invalid params - unable to parse guid")))
        })
    })
})
