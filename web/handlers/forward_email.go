package handlers

import (
    "net/http"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/config"
    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/log"
    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/requests"
    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/uaa"
)

type ForwardEmail struct {
    requestBuilder requests.RequestBuilderInterface
    requestSender  requests.RequestSenderInterface
}

func NewForwardEmail(requestBuilder requests.RequestBuilderInterface,
    requestSender requests.RequestSenderInterface) ForwardEmail {

    return ForwardEmail{
        requestBuilder: requestBuilder,
        requestSender:  requestSender,
    }
}

func (handler ForwardEmail) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    if req.Body != nil {
        params := make(map[string]string)

        log.Println("Sendgrid request header: ", req.Header)
        log.Printf("The body of the post: %#v", req.Body)

        err := req.ParseMultipartForm(8096)
        if err != nil {
            panic(err)
        }

        params["to"] = req.MultipartForm.Value["to"][0]
        params["text"] = "Eventually the text from the sendgrid post..."

        env := config.NewEnvironment()
        uaa := uaa.NewUAAClient(env)

        request, err := handler.requestBuilder.Build(params, uaa.AccessToken())

        if err != nil {
            log.PrintlnErr("Panicking with error: " + err.Error())
            panic(err) // TODO HANDLE THE ERROR CORRECTLY
        }
        handler.requestSender.Send(request)
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{}`))
}
