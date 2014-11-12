package requests_test

import (
	"errors"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/requests"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestSender", func() {

	var sender requests.RequestSender
	var request *http.Request
	var err error

	BeforeEach(func() {
		request, err = http.NewRequest("get", "http://example.com/testing", nil)
		if err != nil {
			panic(err)
		}

		sender = requests.NewRequestSender()
	})

	Context("when the request to notifications is successful", func() {
		BeforeEach(func() {
			sender.MakeRequest = func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
				}, nil
			}
		})

		It("does not return an error", func() {
			err := sender.Send(request)

			Expect(err).To(BeNil())
		})
	})

	Context("when the request to notifications returns an error", func() {
		Context("when the request returns an error", func() {
			BeforeEach(func() {
				sender.MakeRequest = func(req *http.Request) (*http.Response, error) {
					return &http.Response{}, errors.New("the request failed yo")
				}
			})

			It("returns a NotificationRequestFailed error", func() {
				err := sender.Send(request)

				Expect(err).To(BeAssignableToTypeOf(requests.NotificationRequestFailed{}))
			})
		})

		Context("when the notifications responds with any other non successful error", func() {
			BeforeEach(func() {
				sender.MakeRequest = func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 500,
					}, nil
				}
			})

			It("returns a NotificationRequestFailed error", func() {
				err := sender.Send(request)

				Expect(err).To(BeAssignableToTypeOf(requests.NotificationRequestFailed{}))
				Expect(err.Error()).To(Equal("Request to notifications failed with status code: 500"))
			})
		})
	})
})
