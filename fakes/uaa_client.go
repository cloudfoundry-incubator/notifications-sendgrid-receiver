package fakes

import (
	"errors"

	"github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type UAAClient struct {
	ErrorAlways bool
}

func NewUAAClient() *UAAClient {
	return &UAAClient{}
}

func (fake *UAAClient) GetClientToken() (uaa.Token, error) {
	if fake.ErrorAlways {
		return uaa.Token{}, errors.New("you done goofed")
	}
	return uaa.Token{Access: "the-access-token"}, nil
}
