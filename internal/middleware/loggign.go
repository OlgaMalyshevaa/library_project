package middleware

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func SetupLogger() {
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Logger.SetOutput(os.Stdout)
}
