package log

import (
	"log"
	"os"

	"github.com/cloudfoundry-incubator/notifications-sendgrid-receiver/config"
)

var Logger *log.Logger
var ErrorLogger *log.Logger

func init() {
	env := config.NewEnvironment()
	if len(env.LogFile) > 0 {
		file, err := os.OpenFile(env.LogFile, os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			panic(err)
		}
		Logger = log.New(file, "[notification-sendgrid-receiver] ", log.LstdFlags)
		ErrorLogger = log.New(file, "[notification-sendgrid-receiver] ", log.LstdFlags)
	} else {
		Logger = log.New(os.Stdout, "\x1b[0m[notification-sendgrid-receiver] ", log.LstdFlags)
		ErrorLogger = log.New(os.Stderr, "\x1b[31m[notification-sendgrid-receiver] ", log.LstdFlags)
	}
}

func Print(v ...interface{}) {
	Logger.Print(v...)
}

func Printf(format string, v ...interface{}) {
	Logger.Printf(format, v...)
}

func Println(v ...interface{}) {
	Logger.Println(v...)
}

func PrintlnErr(v ...interface{}) {
	ErrorLogger.Println(v...)
}
