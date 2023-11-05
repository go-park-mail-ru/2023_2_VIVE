package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger(logFile *os.File) *logrus.Logger {
	jsonFormatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}

	logger := logrus.New()
	logger.SetOutput(logFile)
	logger.SetFormatter(jsonFormatter)

	return logger
}
