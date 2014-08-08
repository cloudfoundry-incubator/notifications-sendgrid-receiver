package requests

import (
    "errors"
    "net/http"
)

type RequestBodyParserInterface interface {
    Parse(*http.Request) (map[string]string, error)
}

type RequestBodyParser struct{}

func (parser RequestBodyParser) Parse(req *http.Request) (map[string]string, error) {
    params := make(map[string]string)
    err := req.ParseMultipartForm(8096)
    if err != nil {
        return nil, errors.New("Could not parse the request body as a form data: " + err.Error())
    }

    if len(req.MultipartForm.Value["to"]) == 0 {
        return nil, errors.New("Could not parse a to field out of the form data")
    }
    params["to"] = req.MultipartForm.Value["to"][0]
    params["text"] = "Eventually the text from the sendgrid post..."
    return params, nil
}
