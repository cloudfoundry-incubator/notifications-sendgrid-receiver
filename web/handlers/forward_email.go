package handlers

import (
    "encoding/json"
    "io/ioutil"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/config"
    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/uaa"
)

type ForwardEmail struct {
    Api SpaceMailerAPIInterface
}

func NewForwardEmail(api SpaceMailerAPIInterface) ForwardEmail {
    return ForwardEmail{
        Api: api,
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
        handler.Api.PostToSpace(uaa.AccessToken(), params)
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{}`))
}
