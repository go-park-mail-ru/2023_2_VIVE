package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func InitLogger(logFile *os.File) {
	jsonFormatter := &logrus.JSONFormatter{
		TimestampFormat: "Mon, 02 Jan 2006 15:04:05 MST",
	}

	// logger := logrus.New()
	Logger.SetOutput(logFile)
	Logger.SetFormatter(jsonFormatter)
	Logger.SetLevel(logrus.DebugLevel)
	Logger = Logger.WithFields(logrus.Fields{"app_name": "app"}).Logger

	// return Logger
}
