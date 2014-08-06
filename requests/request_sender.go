package requests

import (
    "fmt"
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
    if err != nil {
        //TODO: eliminate this print
        log.PrintlnErr("Request to Notification server failed: " + err.Error())
        return NewNotificationRequestFailed(err.Error())
    }

    log.Printf("notifications response code: %d", response.StatusCode)

    if response.StatusCode != 200 {
        return NewNotificationRequestFailed(fmt.Sprintf("Request to notifications failed with status code: %d", response.StatusCode))
    }

    return nil
}
