package services_test

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Request Builder", func() {
	var request *http.Request
	var err error
	var builder services.RequestBuilder
	var params services.RequestParams
	var notificationsHost string

	BeforeEach(func() {
		params = services.RequestParams{
			Kind:    "bananapanic.com",
			From:    "thisperson@example.com",
			ReplyTo: "dog@example.com",
			To:      "space-guid-mammoth1-banana2-damage3@example.com",
			Subject: "Hamhock",
			Text:    "the text of the email",
			HTML:    "the html of the email",
		}

		notificationsHost = "notifications.example.com"
		builder = services.NewRequestBuilder(notificationsHost)

		request, err = builder.Build(params.ToMap(), "the-access-token")
		if err != nil {
			panic(err)
		}
	})

	Context("when the space guid cannot be parsed", func() {
		It("returns an error", func() {
			params.To = "fake-banHammer-____&&&989867.com"
			_, err := builder.Build(params.ToMap(), "the-access-token")
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("Invalid params - unable to parse guid"))
		})
	})

	Context("when the space guid can be parsed", func() {
		It("returns a post request to the appropriate endpoint", func() {
			Expect(request.Method).To(Equal("POST"))
			Expect(request.URL.String()).To(Equal(notificationsHost + "/spaces/mammoth1-banana2-damage3"))
		})
	})

	It("returns a request with the appropriate headers", func() {
		Expect(request.Header.Get("Authorization")).To(Equal("Bearer the-access-token"))
		Expect(request.Header.Get("Content-Type")).To(Equal("application/json"))
	})

	It("returns a request with the appropriate json body", func() {
		buffer := bytes.NewBuffer([]byte{})
		buffer.ReadFrom(request.Body)

		jsonBody := make(map[string]string)
		err := json.Unmarshal(buffer.Bytes(), &jsonBody)
		if err != nil {
			panic(err)
		}

		Expect(jsonBody["kind_id"]).To(Equal("bananapanic.com"))
		Expect(jsonBody["reply_to"]).To(Equal("dog@example.com"))
		Expect(jsonBody["subject"]).To(Equal("Hamhock"))
		Expect(jsonBody["text"]).To(Equal("the text of the email"))
		Expect(jsonBody["html"]).To(Equal("the html of the email"))
	})
})
