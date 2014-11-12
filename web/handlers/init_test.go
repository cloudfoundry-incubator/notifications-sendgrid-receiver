package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/requests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebHandlersSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Web Handlers Suite")
}

func FakeOAuth() (token string) {
	return "fakeTokenThatNeedsToBeCreatedOrSomething"
}

type FakeRequestBuilder struct {
	Request     *http.Request
	Params      map[string]string
	ErrorAlways bool
}

func (fake *FakeRequestBuilder) Build(params requests.RequestParams, accessToken string) (*http.Request, error) {
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

type FakeRequestSender struct {
	Request *http.Request
}

func (fake *FakeRequestSender) Send(request *http.Request) error {
	fake.Request = request
	return nil
}

type FakeUAAClient struct {
	ErrorAlways bool
}

func (fake *FakeUAAClient) AccessToken() (string, error) {
	if fake.ErrorAlways {
		return "", errors.New("you done goofed")
	}
	return "the-access-token", nil
}

type FakeRequestBodyParser struct {
	ErrorAlways bool
	Params      requests.RequestParams
}

func (fake *FakeRequestBodyParser) Parse(req *http.Request) (requests.RequestParams, error) {
	if fake.ErrorAlways {
		return fake.Params, errors.New("an error occured")
	}
	return fake.Params, nil
}

type FakeBasicAuthenticator struct {
	InvalidAuth bool
}

func (fake FakeBasicAuthenticator) Verify(header http.Header) bool {
	if fake.InvalidAuth {
		return false
	}
	return true
}
