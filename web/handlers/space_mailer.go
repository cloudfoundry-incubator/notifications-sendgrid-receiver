package handlers

import (
    "bytes"
    "net/http"
    "net/url"
    "os"
    "strconv"
)

type SpaceMailerAPIInterface interface {
    PostToSpace(string, map[string]string) error
}

type SpaceMailerAPI struct{}

func (api SpaceMailerAPI) PostToSpace(uaaAccessToken string, params map[string]string) error {
    host := os.Getenv("NOTIFICATIONS_SERVER")

    spaceGuid, err := parseSpaceGuid(params["to"])
    if err != nil {
        return err
    }

    spaceURL := host + "/spaces/" + spaceGuid

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
