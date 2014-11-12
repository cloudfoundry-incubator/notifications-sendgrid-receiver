package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/log"
	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/requests"
	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/uaa"
)

type ForwardEmail struct {
	requestBuilder     requests.RequestBuilderInterface
	requestSender      requests.RequestSenderInterface
	uaaClient          uaa.UAAClientInterface
	requestBodyParser  requests.RequestBodyParserInterface
	basicAuthenticator requests.BasicAuthenticatorInterface
}

func NewForwardEmail(requestBuilder requests.RequestBuilderInterface,
	requestSender requests.RequestSenderInterface,
	uaa uaa.UAAClientInterface,
	requestBodyParser requests.RequestBodyParserInterface,
	basicAuthenticator requests.BasicAuthenticatorInterface) ForwardEmail {

	return ForwardEmail{
		requestBuilder:     requestBuilder,
		requestSender:      requestSender,
		uaaClient:          uaa,
		requestBodyParser:  requestBodyParser,
		basicAuthenticator: basicAuthenticator,
	}
}

func (handler ForwardEmail) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{}`))
		return
	}

	authenticates := handler.basicAuthenticator.Verify(req.Header)
	if authenticates == false {
		log.PrintlnErr("401 StatusUnauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params, err := handler.requestBodyParser.Parse(req)
	if err != nil {
		log.PrintlnErr(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}

	accessToken, err := handler.uaaClient.AccessToken()
	if err != nil {
		log.PrintlnErr("UAA returned an error: " + err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	request, err := handler.requestBuilder.Build(params, accessToken)

	if err != nil {
		log.PrintlnErr("Build request failed with error: " + err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	handler.requestSender.Send(request)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
