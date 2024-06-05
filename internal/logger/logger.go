package logger

import (
	"github.com/sirupsen/logrus"

	"github.com/KillReall666/Rutube-project/internal/config"
)

type Logger struct {
	log *logrus.Logger
	cfg *config.Config
}

func New() *Logger {
	log := logrus.New()
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	return &Logger{
		log: log,
	}
}

func (l *Logger) LogInfo(message string, args ...interface{}) {
	l.log.Infoln(message, args)
}

func (l *Logger) LogError(message string, args ...interface{}) {
	l.log.Error(message, args)
}

func (l *Logger) LogFatal(message string, args ...interface{}) {
	l.log.Fatalln(message, args)
}
