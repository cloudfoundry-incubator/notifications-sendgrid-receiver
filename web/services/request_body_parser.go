package services

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
)

type RequestBodyParser struct{}

func NewRequestBodyParser() RequestBodyParser {
	return RequestBodyParser{}
}

func (parser RequestBodyParser) Parse(req *http.Request) (map[string]string, error) {
	params := RequestParams{}

	err := req.ParseMultipartForm(8096)
	if err != nil {
		return map[string]string{}, errors.New("Could not parse the request body as a form data: " + err.Error())
	}

	if len(req.MultipartForm.Value["to"]) == 0 {
		return map[string]string{}, errors.New("Could not parse a to field out of the form data")
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

	params.ReplyTo = params.From

	if len(req.MultipartForm.Value["headers"]) != 0 {
		headers := strings.Split(req.MultipartForm.Value["headers"][0], "\n")
		for _, value := range headers {
			if strings.HasPrefix(value, "Reply-To: ") {
				params.ReplyTo = strings.TrimPrefix(value, "Reply-To: ")
			}
		}
	}

	regex := regexp.MustCompile("@(.*[^>])")
	if len(regex.FindStringSubmatch(params.ReplyTo)) > 0 {
		params.Kind = regex.FindStringSubmatch(params.ReplyTo)[1]
	}

	return params.ToMap(), nil
}
