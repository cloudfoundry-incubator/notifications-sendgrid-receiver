package handlers

import (
    "bytes"
    "errors"
    "net/http"
    "net/url"
    "regexp"
    "strconv"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/config"
)

type NotificationsReceiverInterface interface {
    PostToSpace(string, map[string]string) error
}

type NotificationsReceiver struct{}

func (receiver NotificationsReceiver) PostToSpace(uaaAccessToken string, params map[string]string) error {
    env := config.NewEnvironment()

    spaceGuid, err := receiver.parseSpaceGuid(params["to"])
    if err != nil {
        return err
    }

    spaceURL := env.NotificationsHost + "/spaces/" + spaceGuid

    data := url.Values{}
    data.Set("kind", params["kind"])
    data.Add("text", params["text"])

    request, err := http.NewRequest("POST", spaceURL, bytes.NewBufferString(data.Encode()))
    if err != nil {
        return err
    }

    request.Header.Set("Authorization", "Bearer "+uaaAccessToken)
    request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return err
    }

    response.Body.Close()
    return nil
}

func (receiver NotificationsReceiver) parseSpaceGuid(email string) (guid string, err error) {
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
