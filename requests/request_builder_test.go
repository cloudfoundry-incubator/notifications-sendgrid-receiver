package requests_test

import (
    "bytes"
    "encoding/json"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/config"
    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/requests"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Request Builder", func() {
    var request *http.Request
    var err error
    var builder requests.RequestBuilder
    var params map[string]string
    var env config.Environment

    BeforeEach(func() {
        params = make(map[string]string)
        params["to"] = "space-guid-mammoth1-banana2-damage3@example.com"
        params["text"] = "the text of the email"
        env = config.NewEnvironment()

        request, err = builder.Build(params, "the-access-token")
        if err != nil {
            panic(err)
        }
    })

    Context("when the space guid cannot be parsed", func() {
        It("returns an error", func() {
            params["to"] = "fake-banHammer-____&&&989867.com"
            _, err := builder.Build(params, "the-access-token")
            Expect(err).ToNot(BeNil())
            Expect(err.Error()).To(Equal("Invalid params - unable to parse guid"))
        })
    })

    Context("when the space guid can be parsed", func() {
        It("returns a post request to the appropriate endpoint", func() {
            Expect(request.Method).To(Equal("POST"))
            Expect(request.URL.String()).To(Equal(env.NotificationsHost + "/spaces/mammoth1-banana2-damage3"))
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

        Expect(jsonBody["text"]).To(Equal("the text of the email"))
        Expect(jsonBody["kind"]).To(Equal("sendgrid-kind-value"))
    })
})
