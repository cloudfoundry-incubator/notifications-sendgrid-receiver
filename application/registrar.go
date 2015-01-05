package application

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Registrar struct {
	notificationsDomain string
}

func NewRegistrar(notificationsDomain string) Registrar {
	return Registrar{
		notificationsDomain: notificationsDomain,
	}
}

func (r Registrar) Register(token string) error {
	var body struct {
		SourceName string `json:"source_name"`
	}
	body.SourceName = "3rd Party Services"
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	statusCode, _, err := makeRequest("PUT", fmt.Sprintf("%s/notifications", r.notificationsDomain), bytes.NewBuffer(bodyBytes), token)
	if err != nil {
		return err
	}
	if statusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Unexpected HTTP status code when registering notification: %d", statusCode))
	}
	return nil
}

func makeRequest(method, path string, content io.Reader, token string) (int, io.Reader, error) {
	request, err := http.NewRequest(method, path, content)
	if err != nil {
		return 0, nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, nil, err
	}

	return response.StatusCode, response.Body, nil
}
