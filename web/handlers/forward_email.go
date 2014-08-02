package handlers

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "net/http"
    "os"
    "regexp"

    "github.com/pivotal-cf/uaa-sso-golang/uaa"
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

        uaaAccessToken := GetUAAToken()
        handler.Api.PostToSpace(uaaAccessToken, params)
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{}`))
}

func parseSpaceGuid(email string) (guid string, err error) {
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

func GetUAAToken() string {
    loginURL := ""
    uaaHost := os.Getenv("UAA_HOST")
    uaaClientID := os.Getenv("UAA_CLIENT_ID")
    uaaClientSecret := os.Getenv("UAA_CLIENT_SECRET")
    accessToken := ""

    uaa := uaa.NewUAA(loginURL, uaaHost, uaaClientID, uaaClientSecret, accessToken)

    uaaToken, err := uaa.GetClientToken()
    if err != nil {
        panic(err)
    }

    return uaaToken.Access
}
