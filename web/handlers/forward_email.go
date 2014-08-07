package handlers

import (
    "bytes"
    "fmt"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/log"
    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/requests"
    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/uaa"
)

type ForwardEmail struct {
    requestBuilder requests.RequestBuilderInterface
    requestSender  requests.RequestSenderInterface
    uaaClient      uaa.UAAClientInterface
}

func NewForwardEmail(requestBuilder requests.RequestBuilderInterface,
    requestSender requests.RequestSenderInterface,
    uaa uaa.UAAClientInterface) ForwardEmail {

    return ForwardEmail{
        requestBuilder: requestBuilder,
        requestSender:  requestSender,
        uaaClient:      uaa,
    }
}

func (handler ForwardEmail) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    if req.Body == nil {

        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`{}`))
        return
    }

    params := make(map[string]string)
    log.Printf("Request Headers: ", req.Header)

    bodyBuffer := &bytes.Buffer{}
    bodyBuffer.ReadFrom(req.Body)
    log.Printf("Request body: %#v", bodyBuffer.String())

    err := req.ParseMultipartForm(8096)
    if err != nil {
        fmt.Println(req.MultipartForm)
        log.PrintlnErr("Could not parse the request body as a form data: " + err.Error())
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if len(req.MultipartForm.Value["to"]) == 0 {
        log.PrintlnErr("Could not parse a to field out of the form data")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    params["to"] = req.MultipartForm.Value["to"][0]
    params["text"] = "Eventually the text from the sendgrid post..."

    accessToken, err := handler.uaaClient.AccessToken()
    if err != nil {
        log.PrintlnErr("UAA returned an error: " + err.Error())
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }

    request, err := handler.requestBuilder.Build(params, accessToken)

    if err != nil {
        log.PrintlnErr("Build request failed with error: " + err.Error())
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    handler.requestSender.Send(request)

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{}`))
}
