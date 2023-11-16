package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger(logFile *os.File) *logrus.Logger {
	jsonFormatter := &logrus.JSONFormatter{
		TimestampFormat: "Mon, 02 Jan 2006 15:04:05 MST",
	}

	logger := logrus.New()
	logger.SetOutput(logFile)
	logger.SetFormatter(jsonFormatter)

	return logger
}
