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
    var validPostBeginning, postSubject, postFrom, postReplyTo, postText, postHTML, validPostEnding string

    BeforeEach(func() {
        validPostBeginning = "--xYzZy\nContent-Disposition: form-data; name=\"to\"\n\nspace-guid-the-guid-88@bananahamhock.com\n--xYzZy"
        postSubject = "\nContent-Disposition: form-data; name=\"subject\"\n\nThis is a great subject\n--xYzZy"
        postFrom = "\nContent-Disposition: form-data; name=\"from\"\n\nincoming-from@example.com\n--xYzZy"
        postReplyTo = "\nContent-Disposition: form-data; name=\"reply-to\"\n\nincoming-reply-to@example.com\n--xYzZy"
        postText = "\nContent-Disposition: form-data; name=\"text\"\n\nThis is the text of the email or something\n--xYzZy"
        postHTML = "\nContent-Disposition: form-data; name=\"html\"\n\n<h1>This is the html of the email</h1>\n--xYzZy"
        validPostEnding = "--\n\n"
    })

    It("sets the RequestParams from the incoming form-data request", func() {
        validPostBody = validPostBeginning + postSubject + postFrom + postReplyTo + postText + postHTML + validPostEnding

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

        Expect(params.To).To(Equal("space-guid-the-guid-88@bananahamhock.com"))
        Expect(params.Subject).To(Equal("This is a great subject"))
        Expect(params.From).To(Equal("incoming-from@example.com"))
        Expect(params.ReplyTo).To(Equal("incoming-reply-to@example.com"))
        Expect(params.Text).To(Equal("This is the text of the email or something"))
        Expect(params.HTML).To(Equal("<h1>This is the html of the email</h1>"))
    })

    It("does not error when the post is missing non-required fields", func() {
        validPostBody = validPostBeginning + validPostEnding

        body := []byte(validPostBody)
        request, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
        if err != nil {
            panic(err)
        }

        request.Header.Add("Content-Type", "multipart/form-data; boundary=xYzZy")

        _, err = bodyParser.Parse(request)
        Expect(err).To(BeNil())
    })

    It("sets the reply-to to be the from when no reply-to is specified", func() {
        validPostBody = validPostBeginning + postSubject + postFrom + postText + validPostEnding

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

        Expect(params.ReplyTo).To(Equal("incoming-from@example.com"))
    })

    It("sets the kind_id to the domain in reply-to (or from) email", func() {
        validPostBody = validPostBeginning + postSubject + postFrom + postReplyTo + postText + validPostEnding

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

        Expect(params.Kind).To(Equal("example.com"))
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
