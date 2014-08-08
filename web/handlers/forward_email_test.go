package handlers_test

import (
    "bytes"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Forward", func() {
    var handler handlers.ForwardEmail
    var fakeRequestBuilder FakeRequestBuilder
    var fakeRequestSender FakeRequestSender
    var fakeUAA FakeUAAClient
    var fakeBodyParser FakeRequestBodyParser

    BeforeEach(func() {
        handler = handlers.NewForwardEmail(&fakeRequestBuilder, &fakeRequestSender, &fakeUAA, &fakeBodyParser)
    })

    AfterEach(func() {
        fakeBodyParser.ErrorAlways = false
        fakeRequestBuilder.ErrorAlways = false
    })

    Describe("ServeHTTP", func() {
        var formData string

        BeforeEach(func() {
            formData = "--xYzZy\nContent-Disposition: form-data; name=\"to\"\n\nspace-guid-the-guid-88@bananahamhock.com\n--xYzZy--\n"
            fakeBodyParser.Params = make(map[string]string)
            fakeBodyParser.Params["to"] = "space-guid-the-guid-88@bananahamhock.com"
        })

        It("sends a request built by the request builder to the notifications service", func() {
            writer := httptest.NewRecorder()
            body := []byte(formData)
            request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
            if err != nil {
                panic(err)
            }

            request.Header.Add("Content-Type", "multipart/form-data; boundary=xYzZy")

            handler.ServeHTTP(writer, request)

            Expect(fakeRequestBuilder.Params["to"]).To(Equal("space-guid-the-guid-88@bananahamhock.com"))
            Expect(fakeRequestSender.Request).To(Equal(fakeRequestBuilder.Request))
        })

        Context("when notifications responds with a success", func() {
            It("returns a 200 response code and an empty JSON body", func() {
                writer := httptest.NewRecorder()
                body := []byte(formData)
                request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
                if err != nil {
                    panic(err)
                }

                request.Header.Add("Content-Type", "multipart/form-data; boundary=xYzZy")
                handler.ServeHTTP(writer, request)

                Expect(writer.Code).To(Equal(http.StatusOK))
                Expect(writer.Body.String()).To(Equal("{}"))
            })
        })

        Context("when no post body is passed", func() {
            It("sets the status code to 400", func() {
                writer := httptest.NewRecorder()
                request, err := http.NewRequest("POST", "/", nil)
                if err != nil {
                    panic(err)
                }

                handler.ServeHTTP(writer, request)

                Expect(writer.Code).To(Equal(http.StatusBadRequest))
            })
        })

        Context("when the body parser returns an error", func() {
            It("sets the status code to 400", func() {
                fakeBodyParser.ErrorAlways = true

                writer := httptest.NewRecorder()
                body := []byte(formData)
                request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
                if err != nil {
                    panic(err)
                }

                handler.ServeHTTP(writer, request)

                Expect(writer.Code).To(Equal(http.StatusBadRequest))
            })
        })

        Context("when the request builder responds with an error", func() {
            It("sets the status code to 503", func() {
                fakeRequestBuilder.ErrorAlways = true

                writer := httptest.NewRecorder()
                body := []byte(formData)
                request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
                if err != nil {
                    panic(err)
                }

                request.Header.Add("Content-Type", "multipart/form-data; boundary=xYzZy")
                handler.ServeHTTP(writer, request)

                Expect(writer.Code).To(Equal(http.StatusServiceUnavailable))
            })
        })

        Context("when uaa is down", func() {
            It("sets the status code to 503", func() {
                fakeUAA.ErrorAlways = true

                writer := httptest.NewRecorder()
                body := []byte(formData)
                request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
                if err != nil {
                    panic(err)
                }

                request.Header.Add("Content-Type", "multipart/form-data; boundary=xYzZy")
                handler.ServeHTTP(writer, request)

                Expect(writer.Code).To(Equal(http.StatusServiceUnavailable))
            })
        })
    })
})
