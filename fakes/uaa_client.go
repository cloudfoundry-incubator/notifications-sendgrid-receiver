package fakes

import "errors"

type UAAClient struct {
	ErrorAlways bool
}

func NewUAAClient() *UAAClient {
	return &UAAClient{}
}

func (fake *UAAClient) AccessToken() (string, error) {
	if fake.ErrorAlways {
		return "", errors.New("you done goofed")
	}
	return "the-access-token", nil
}
