package application

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Registrar struct {
	notificationsDomain string
	httpClient          *http.Client
}

func NewRegistrar(notificationsDomain string, verifySSL bool) Registrar {
	return Registrar{
		notificationsDomain: notificationsDomain,
		httpClient:          GetClient(verifySSL),
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
	statusCode, _, err := r.makeRequest("PUT", fmt.Sprintf("%s/notifications", r.notificationsDomain), bytes.NewBuffer(bodyBytes), token)
	if err != nil {
		return err
	}
	if statusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("Unexpected HTTP status code when registering notification: %d", statusCode))
	}
	return nil
}

func (r Registrar) makeRequest(method, path string, content io.Reader, token string) (int, io.Reader, error) {
	request, err := http.NewRequest(method, path, content)
	if err != nil {
		return 0, nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	response, err := r.httpClient.Do(request)
	if err != nil {
		return 0, nil, err
	}

	return response.StatusCode, response.Body, nil
}

var _client *http.Client
var mutex sync.Mutex

func GetClient(verifySSL bool) *http.Client {
	mutex.Lock()
	defer mutex.Unlock()

	if _client == nil {
		_client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: !verifySSL,
				},
			},
		}
	}

	return _client
}
