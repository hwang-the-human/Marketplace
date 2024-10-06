package config

import (
	"github.com/sirupsen/logrus"
	"os"
)

func InitLogrus() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logrus.SetLevel(logrus.InfoLevel)

	logrus.SetOutput(os.Stdout)
}
