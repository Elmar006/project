package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func Init() {
	log.SetOutput(os.Stdout)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   false,
	})

	log.SetLevel(log.InfoLevel)
}

func L() *log.Logger {
	return log.StandardLogger()
}
