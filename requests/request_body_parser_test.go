package requests_test

import (
    "bytes"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/requests"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("RequestBodyParser", func() {
    var validPostBody string
    var bodyParser requests.RequestBodyParser

    BeforeEach(func() {
        validPostBody = "--xYzZy\nContent-Disposition: form-data; name=\"to\"\n\nspace-guid-the-guid-88@bananahamhock.com\n--xYzZy--\n"
    })

    It("sets the correct values in params", func() {
        body := []byte(validPostBody)
        request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
        if err != nil {
            panic(err)
        }

        request.Header.Add("Content-Type", "multipart/form-data; boundary=xYzZy")

        params, err := bodyParser.Parse(request)
        if err != nil {
            panic(err)
        }
        Expect(params["to"]).To(Equal("space-guid-the-guid-88@bananahamhock.com"))
        Expect(params["text"]).To(Equal("Eventually the text from the sendgrid post..."))
    })

    Context("when the to parameter is missing from the post body", func() {
        It("return an error", func() {
            body := []byte("--xYzZy\nContent-Disposition: form-data; name=\"not-to\"\n\nspace-guid-the-guid-88@bananahamhock.com\n--xYzZy--\n")
            request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
            if err != nil {
                panic(err)
            }

            request.Header.Add("Content-Type", "multipart/form-data; boundary=xYzZy")

            _, err = bodyParser.Parse(request)
            Expect(err).ToNot(BeNil())
        })
    })

    Context("when the multipart data cannot be parsed", func() {
        It("returns an error", func() {
            body := []byte("not a valid multipart")
            request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
            if err != nil {
                panic(err)
            }

            request.Header.Add("Content-Type", "multipart/form-data; boundary=xYzZy")

            _, err = bodyParser.Parse(request)
            Expect(err).ToNot(BeNil())
        })
    })
})
