package requests

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "regexp"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/config"
)

type RequestBuilderInterface interface {
    Build(RequestParams, string) (*http.Request, error)
}

type RequestBuilder struct{}

func NewRequestBuilder() RequestBuilder {
    return RequestBuilder{}
}

func (builder RequestBuilder) Build(params RequestParams, accessToken string) (*http.Request, error) {
    guid, err := builder.parseSpaceGuid(params.To)
    if err != nil {
        return &http.Request{}, err
    }

    env := config.NewEnvironment()
    notificationEndpoint := env.NotificationsHost + "/spaces/" + guid

    body := make(map[string]string)
    body["kind"] = params.Kind
    body["from"] = params.From
    body["replyTo"] = params.ReplyTo
    body["to"] = params.To
    body["subject"] = params.Subject
    body["text"] = params.Text
    body["html"] = params.HTML

    jsonBody, err := json.Marshal(body)
    if err != nil {
        return &http.Request{}, err
    }

    request, err := http.NewRequest("POST", notificationEndpoint, bytes.NewBufferString(string(jsonBody)))
    if err != nil {
        return &http.Request{}, err
    }

    request.Header.Set("Authorization", "Bearer "+accessToken)
    request.Header.Set("Content-Type", "application/json")
    return request, nil
}

func (builder RequestBuilder) parseSpaceGuid(email string) (guid string, err error) {
    regex, err := regexp.Compile("space-guid-([a-zA-Z0-9-]*)@")
    if err != nil {
        return "", err
    }

    if regex.FindStringSubmatch(email) == nil {
        return "", errors.New("Invalid params - unable to parse guid")
    }

    guid = regex.FindStringSubmatch(email)[1]
    return
}
