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
	BasicAuthUserName string
	BasicAuthPassword string
	CCHost            string
	Port              string
	VerifySSL         bool
	NotificationsHost string
	LogFile           string
}

func NewEnvironment() Environment {
	return Environment{
		UAAHost:           loadOrPanic("UAA_HOST"),
		UAAClientID:       loadOrPanic("UAA_CLIENT_ID"),
		UAAClientSecret:   loadOrPanic("UAA_CLIENT_SECRET"),
		BasicAuthUserName: loadOrPanic("BASIC_AUTH_USER_NAME"),
		BasicAuthPassword: loadOrPanic("BASIC_AUTH_PASSWORD"),
		Port:              loadOrPanic("PORT"),
		CCHost:            loadOrPanic("CC_HOST"),
		VerifySSL:         loadBool("VERIFY_SSL", true),
		NotificationsHost: loadOrPanic("NOTIFICATIONS_HOST"),
		LogFile:           loadOrDefault("LOG_FILE", ""),
	}
}

func loadOrPanic(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic(errors.New(fmt.Sprintf("Could not find required %s environment variable", name)))
	}
	return value
}

func loadOrDefault(name, defaultValue string) string {
	variable := os.Getenv(name)
	if variable == "" {
		variable = defaultValue
	}
	return variable
}

func loadBool(name string, defaultValue bool) bool {
	value, err := strconv.ParseBool(os.Getenv(name))
	if err != nil {
		return defaultValue
	}

	return value
}
