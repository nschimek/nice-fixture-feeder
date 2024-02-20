package core

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = initializeLogger()

func initializeLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})

	file, err := os.OpenFile("out/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err == nil {
		log.SetOutput(file)
	}

	return log
}
