package handlers

import (
    "encoding/json"
    "io/ioutil"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/config"
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
        var body []byte
        var params map[string]string

        body, _ = ioutil.ReadAll(req.Body)
        json.Unmarshal(body, &params)

        env := config.NewEnvironment()
        uaa := uaa.NewUAAClient(env)

        request, err := handler.requestBuilder.Build(params, uaa.AccessToken())

        if err != nil {
            panic(err) // TODO HANDLE THE ERROR CORRECTLY
        }
        handler.requestSender.Send(request)
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{}`))
}
