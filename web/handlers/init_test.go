package handlers_test

import (
    "bytes"
    "errors"
    "net/http"
    "testing"

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

func (fake *FakeRequestBuilder) Build(params map[string]string, accessToken string) (*http.Request, error) {
    if fake.ErrorAlways {
        return &http.Request{}, errors.New("Fake Request Builder Error")
    }
    fake.Params = params
    request, err := http.NewRequest("POST", "the host", bytes.NewBufferString(params["body"]))
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
