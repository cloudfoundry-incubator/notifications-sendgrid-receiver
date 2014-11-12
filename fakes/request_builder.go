package fakes

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/requests"
)

type RequestBuilder struct {
	Request     *http.Request
	Params      map[string]string
	ErrorAlways bool
}

func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{}
}

func (fake *RequestBuilder) Build(params requests.RequestParams, accessToken string) (*http.Request, error) {
	if fake.ErrorAlways {
		return &http.Request{}, errors.New("Fake Request Builder Error")
	}
	fake.Params = map[string]string{}

	fake.Params["to"] = params.To

	jsonBody, err := json.Marshal(fake.Params)
	if err != nil {
		return &http.Request{}, err
	}

	request, err := http.NewRequest("POST", "the host", bytes.NewBufferString(string(jsonBody)))
	if err != nil {
		panic(err)
	}
	request.Header.Set("fake Headers", "the fake headers")
	fake.Request = request
	return request, nil
}
