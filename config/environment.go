package config

import (
    "errors"
    "fmt"
    "os"
    "strconv"
)

var UAAPublicKey string

type Environment struct {
    UAAHost           string
    UAAClientID       string
    UAAClientSecret   string
    CCHost            string
    VerifySSL         bool
    RootPath          string
    NotificationsHost string
} 
func NewEnvironment() Environment {
    return Environment{
        UAAHost:           loadOrPanic("UAA_HOST"),
        UAAClientID:       loadOrPanic("UAA_CLIENT_ID"),
        UAAClientSecret:   loadOrPanic("UAA_CLIENT_SECRET"),
        CCHost:            loadOrPanic("CC_HOST"),
        VerifySSL:         loadBool("VERIFY_SSL", true),
        RootPath:          loadOrPanic("ROOT_PATH"),
        NotificationsHost: loadOrPanic("NOTIFICATIONS_HOST"),
    }
}

func loadOrPanic(name string) string {
    value := os.Getenv(name)
    if value == "" {
        panic(errors.New(fmt.Sprintf("Could not find required %s environment variable", name)))
    }
    return value
}

func loadBool(name string, defaultValue bool) bool {
    value, err := strconv.ParseBool(os.Getenv(name))
    if err != nil {
        return defaultValue
    }

    return value
}
