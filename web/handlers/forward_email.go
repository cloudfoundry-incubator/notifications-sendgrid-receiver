package handlers

import (
	"log"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/web/services"
	"github.com/ryanmoran/stack"
)

type ForwardEmail struct {
	requestBuilder     RequestBuilderInterface
	requestSender      RequestSenderInterface
	uaaClient          UAAClientInterface
	requestBodyParser  RequestBodyParserInterface
	basicAuthenticator BasicAuthenticatorInterface
	logger             *log.Logger
}

func NewForwardEmail(requestBuilder RequestBuilderInterface, requestSender RequestSenderInterface,
	uaaClient UAAClientInterface, requestBodyParser RequestBodyParserInterface,
	basicAuthenticator BasicAuthenticatorInterface, logger *log.Logger) ForwardEmail {

	return ForwardEmail{
		requestBuilder:     requestBuilder,
		requestSender:      requestSender,
		uaaClient:          uaaClient,
		requestBodyParser:  requestBodyParser,
		basicAuthenticator: basicAuthenticator,
		logger:             logger,
	}
}

func (handler ForwardEmail) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	if req.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{}`))
		return
	}

	authenticates := handler.basicAuthenticator.Verify(req.Header)
	if authenticates == false {
		handler.logger.Println("401 Unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params, err := handler.requestBodyParser.Parse(req)
	if err != nil {
		handler.logger.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}

	token, err := handler.uaaClient.GetClientToken()
	if err != nil {
		handler.logger.Println("UAA returned an error: " + err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	request, err := handler.requestBuilder.Build(params, token.Access)
	if err != nil {
		handler.logger.Println("Build request failed with error: " + err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	err = handler.requestSender.Send(request)
	if err != nil {
		handler.logger.Println("Failed to send request to notifications: " + err.Error())

		if err == services.SpaceNotFound("") {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
