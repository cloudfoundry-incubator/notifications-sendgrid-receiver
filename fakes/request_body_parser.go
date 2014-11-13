package fakes

import (
	"errors"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web/services"
)

type RequestBodyParser struct {
	ErrorAlways bool
	Params      services.RequestParams
}

func NewRequestBodyParser() *RequestBodyParser {
	return &RequestBodyParser{}
}

func (fake *RequestBodyParser) Parse(req *http.Request) (map[string]string, error) {
	if fake.ErrorAlways {
		return fake.Params.ToMap(), errors.New("an error occured")
	}
	return fake.Params.ToMap(), nil
}
