package fakes

import "net/http"

type RequestSender struct {
	Request *http.Request
}

func NewRequestSender() *RequestSender {
	return &RequestSender{}
}

func (fake *RequestSender) Send(request *http.Request) error {
	fake.Request = request
	return nil
}
