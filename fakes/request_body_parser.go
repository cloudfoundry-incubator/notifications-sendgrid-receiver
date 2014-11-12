package fakes

import (
	"errors"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/requests"
)

type RequestBodyParser struct {
	ErrorAlways bool
	Params      requests.RequestParams
}

func NewRequestBodyParser() *RequestBodyParser {
	return &RequestBodyParser{}
}

func (fake *RequestBodyParser) Parse(req *http.Request) (requests.RequestParams, error) {
	if fake.ErrorAlways {
		return fake.Params, errors.New("an error occured")
	}
	return fake.Params, nil
}
