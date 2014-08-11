package requests

import (
    "errors"
    "net/http"
    "regexp"
)

type RequestBodyParserInterface interface {
    Parse(*http.Request) (RequestParams, error)
}

type RequestBodyParser struct {
}

type RequestParams struct {
    To      string
    Subject string
    From    string
    ReplyTo string
    Text    string
    HTML    string
    Kind    string
}

func (parser RequestBodyParser) Parse(req *http.Request) (RequestParams, error) {
    params := RequestParams{}
    err := req.ParseMultipartForm(8096)
    if err != nil {
        return RequestParams{}, errors.New("Could not parse the request body as a form data: " + err.Error())
    }

    if len(req.MultipartForm.Value["to"]) == 0 {
        return RequestParams{}, errors.New("Could not parse a to field out of the form data")
    }

    params.To = req.MultipartForm.Value["to"][0]

    if len(req.MultipartForm.Value["subject"]) != 0 {
        params.Subject = req.MultipartForm.Value["subject"][0]
    }

    if len(req.MultipartForm.Value["from"]) != 0 {
        params.From = req.MultipartForm.Value["from"][0]
    }

    if len(req.MultipartForm.Value["text"]) != 0 {
        params.Text = req.MultipartForm.Value["text"][0]
    }

    if len(req.MultipartForm.Value["html"]) != 0 {
        params.HTML = req.MultipartForm.Value["html"][0]
    }

    if len(req.MultipartForm.Value["reply-to"]) == 0 {
        params.ReplyTo = params.From
    } else {
        params.ReplyTo = req.MultipartForm.Value["reply-to"][0]
    }

    regex := regexp.MustCompile("@(.*[^>])")
    if len(regex.FindStringSubmatch(params.ReplyTo)) > 0 {
        params.Kind = regex.FindStringSubmatch(params.ReplyTo)[1]
    }
    return params, nil
}
