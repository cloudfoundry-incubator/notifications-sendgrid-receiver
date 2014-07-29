package handlers

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "regexp"
)

type ForwardEmail struct {
    Api SpaceMailerApiInterface
}

func NewForwardEmail(api SpaceMailerApiInterface) ForwardEmail {
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

        handler.Api.PostToSpace(params)
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{}`))
}

func parseSpaceGuid(email string) string {
    re := regexp.MustCompile("space-guid-([a-zA-Z0-9-]*)@")
    return re.FindStringSubmatch(email)[1]
}
