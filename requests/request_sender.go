package requests

import (
    "fmt"
    "io/ioutil"
    "net/http"

    "github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/log"
)

type RequestSenderInterface interface {
    Send(*http.Request) error
}

type NotificationRequestFailed struct {
    message string
}

func NewNotificationRequestFailed(message string) NotificationRequestFailed {
    return NotificationRequestFailed{
        message: message,
    }
}

func (err NotificationRequestFailed) Error() string {
    return err.message
}

type RequestSender struct {
    MakeRequest func(*http.Request) (*http.Response, error)
}

func NewRequestSender() RequestSender {
    return RequestSender{
        MakeRequest: func(req *http.Request) (*http.Response, error) {
            client := http.Client{}
            return client.Do(req)
        },
    }
}

func (sender RequestSender) Send(req *http.Request) error {
    response, err := sender.MakeRequest(req)
    log.Println("Outgoing request to notification-service:", req)
    if err != nil {
        return NewNotificationRequestFailed(err.Error())
    }

    if response.Body != nil {
        errorMessage, _ := ioutil.ReadAll(response.Body)
        log.Printf("notifications response body: %s", string(errorMessage))
    }

    log.Printf("notifications response code: %d", response.StatusCode)

    if response.StatusCode != 200 {
        return NewNotificationRequestFailed(fmt.Sprintf("Request to notifications failed with status code: %d", response.StatusCode))
    }

    return nil
}
