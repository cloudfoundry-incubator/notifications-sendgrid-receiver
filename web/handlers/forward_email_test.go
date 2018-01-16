package handlers_test

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/fakes"
	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web/handlers"
	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Forward", func() {
	var handler handlers.ForwardEmail
	var requestBuilder *fakes.RequestBuilder
	var requestSender *fakes.RequestSender
	var uaaClient *fakes.UAAClient
	var bodyParser *fakes.RequestBodyParser
	var basicAuthenticator *fakes.BasicAuthenticator
	var throttler *fakes.Throttler

	BeforeEach(func() {
		requestBuilder = fakes.NewRequestBuilder()
		requestSender = fakes.NewRequestSender()
		uaaClient = fakes.NewUAAClient()
		bodyParser = fakes.NewRequestBodyParser()
		basicAuthenticator = fakes.NewBasicAuthenticator()
		throttler = &fakes.Throttler{}
		logger := log.New(bytes.NewBuffer([]byte{}), "", 0)

		handler = handlers.NewForwardEmail(requestBuilder, requestSender, uaaClient, bodyParser, basicAuthenticator, throttler,
			logger)
	})

	AfterEach(func() {
		bodyParser.ErrorAlways = false
		requestBuilder.ErrorAlways = false
	})

	Describe("ServeHTTP", func() {
		var formData string

		BeforeEach(func() {
			formData = "--xYzZy\nContent-Disposition: form-data; name=\"to\"\n\nspace-guid-the-guid-88@bananahamhock.com\n--xYzZy--\n"
			bodyParser.Params = services.RequestParams{}
			bodyParser.Params.To = "space-guid-the-guid-88@bananahamhock.com"
		})

		Context("when the basic auth header is invalid", func() {
			BeforeEach(func() {
				basicAuthenticator.InvalidAuth = true
			})

			AfterEach(func() {
				basicAuthenticator.InvalidAuth = false
			})

			It("sets the status code to 401", func() {
				writer := httptest.NewRecorder()
				body := []byte(formData)
				request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}

				handler.ServeHTTP(writer, request, nil)
				Expect(writer.Code).To(Equal(http.StatusUnauthorized))
			})
		})

		It("sends a request built by the request builder to the notifications service", func() {
			writer := httptest.NewRecorder()
			body := []byte(formData)
			request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}

			request.Header.Add("Content-Type", "multipart/form-data; boundary=xYzZy")

			handler.ServeHTTP(writer, request, nil)

			Expect(requestBuilder.Params["to"]).To(Equal("space-guid-the-guid-88@bananahamhock.com"))
			Expect(requestSender.Request).To(Equal(requestBuilder.Request))
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
				handler.ServeHTTP(writer, request, nil)

				Expect(writer.Code).To(Equal(http.StatusOK))
				Expect(writer.Body.String()).To(Equal("{}"))
				Expect(throttler.FinishCalled).To(BeTrue())
			})
		})

		Context("when no post body is passed", func() {
			It("sets the status code to 400", func() {
				writer := httptest.NewRecorder()
				request, err := http.NewRequest("POST", "/", nil)
				if err != nil {
					panic(err)
				}

				handler.ServeHTTP(writer, request, nil)

				Expect(writer.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when the body parser returns an error", func() {
			It("sets the status code to 400", func() {
				bodyParser.ErrorAlways = true

				writer := httptest.NewRecorder()
				body := []byte(formData)
				request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}

				handler.ServeHTTP(writer, request, nil)

				Expect(writer.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when the request builder responds with an error", func() {
			It("sets the status code to 503", func() {
				requestBuilder.ErrorAlways = true

				writer := httptest.NewRecorder()
				body := []byte(formData)
				request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}

				request.Header.Add("Content-Type", "multipart/form-data; boundary=xYzZy")
				handler.ServeHTTP(writer, request, nil)

				Expect(writer.Code).To(Equal(http.StatusServiceUnavailable))
			})
		})

		Context("when the request sender returns an error", func() {
			It("responds to a missing space error with a 200", func() {
				requestSender.SendError = services.SpaceNotFound("this is a failure")

				writer := httptest.NewRecorder()
				body := []byte(formData)
				request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}

				handler.ServeHTTP(writer, request, nil)

				Expect(writer.Code).To(Equal(http.StatusOK))
			})

			It("responds to all other errors with a 503", func() {
				requestSender.SendError = errors.New("There's a snake in my boot!!!")

				writer := httptest.NewRecorder()
				body := []byte(formData)
				request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}

				handler.ServeHTTP(writer, request, nil)

				Expect(writer.Code).To(Equal(http.StatusServiceUnavailable))
			})
		})

		Context("when uaa is down", func() {
			It("sets the status code to 503", func() {
				uaaClient.ErrorAlways = true

				writer := httptest.NewRecorder()
				body := []byte(formData)
				request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}

				request.Header.Add("Content-Type", "multipart/form-data; boundary=xYzZy")
				handler.ServeHTTP(writer, request, nil)

				Expect(writer.Code).To(Equal(http.StatusServiceUnavailable))
			})
		})

		Context("when receiving too many requests", func() {
			It("sets the status code to 429", func() {
				throttler.TooManyReceived = true

				writer := httptest.NewRecorder()
				body := []byte(formData)
				request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
				if err != nil {
					panic(err)
				}

				handler.ServeHTTP(writer, request, nil)

				Expect(writer.Code).To(Equal(http.StatusTooManyRequests))
			})
		})
	})
})
