package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
)

type RequestBuilder struct {
	notificationsHost string
}

func NewRequestBuilder(notificationsHost string) RequestBuilder {
	return RequestBuilder{
		notificationsHost: notificationsHost,
	}
}

func (builder RequestBuilder) Build(paramsMap map[string]string, accessToken string) (*http.Request, error) {
	params := NewRequestParamsFromMap(paramsMap)
	guid, err := builder.parseSpaceGuid(params.To)
	if err != nil {
		return &http.Request{}, err
	}

	notificationEndpoint := builder.notificationsHost + "/spaces/" + guid

	body := make(map[string]string)
	body["kind_id"] = params.Kind
	body["reply_to"] = params.ReplyTo
	body["subject"] = params.Subject
	body["text"] = params.Text
	body["html"] = params.HTML

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return &http.Request{}, err
	}

	request, err := http.NewRequest("POST", notificationEndpoint, bytes.NewBufferString(string(jsonBody)))
	if err != nil {
		return &http.Request{}, err
	}

	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Content-Type", "application/json")
	return request, nil
}

func (builder RequestBuilder) parseSpaceGuid(email string) (guid string, err error) {
	regex, err := regexp.Compile("space-guid-([a-zA-Z0-9-]*)@")
	if err != nil {
		return "", err
	}

	if regex.FindStringSubmatch(email) == nil {
		return "", errors.New("Invalid params - unable to parse guid")
	}

	guid = regex.FindStringSubmatch(email)[1]
	return
}
