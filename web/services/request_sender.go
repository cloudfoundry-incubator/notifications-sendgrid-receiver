package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type NotificationRequestFailed string

func (err NotificationRequestFailed) Error() string {
	return string(err)
}

type RequestSender struct {
	MakeRequest func(*http.Request) (*http.Response, error)
	logger      *log.Logger
}

func NewRequestSender(logger *log.Logger) RequestSender {
	return RequestSender{
		MakeRequest: func(req *http.Request) (*http.Response, error) {
			client := http.DefaultClient
			return client.Do(req)
		},
		logger: logger,
	}
}

func (sender RequestSender) Send(req *http.Request) error {
	response, err := sender.MakeRequest(req)
	sender.logger.Println("Outgoing request to notification-service:", req)
	if err != nil {
		return NotificationRequestFailed(err.Error())
	}

	if response.Body != nil {
		errorMessage, _ := ioutil.ReadAll(response.Body)
		sender.logger.Printf("notifications response body: %s", string(errorMessage))
	}

	sender.logger.Printf("notifications response code: %d", response.StatusCode)

	if response.StatusCode != 200 {
		return NotificationRequestFailed(fmt.Sprintf("Request to notifications failed with status code: %d", response.StatusCode))
	}

	return nil
}
