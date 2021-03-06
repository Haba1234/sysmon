package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger структура логгера.
type Logger struct {
	fl logrus.FieldLogger
}

// New конструктор.
func New(logLevel string, path string) (*Logger, error) {
	log := logrus.New()

	log.Formatter = &logrus.TextFormatter{
		ForceColors:     true,
		DisableColors:   false,
		TimestampFormat: "02/Jan/2006 15:04:05 -0700",
		FullTimestamp:   true,
	}

	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("logger. Error in settings (log file): %w", err)
	}
	writer := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(writer)

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, fmt.Errorf("logger. Error in settings (level): %w", err)
	}
	log.SetLevel(level)

	log.WithFields(logrus.Fields{
		"package": "logger",
		"level":   level,
		"logFile": path,
	}).Debug("Logger setup successful")

	return &Logger{log}, nil
}

// Info информационные сообщения.
func (l Logger) Info(msg string) {
	l.fl.Info(msg)
}

// Warn предупреждающие сообщения.
func (l Logger) Warn(msg string) {
	l.fl.Warn(msg)
}

// Error сообщения об ошибках.
func (l Logger) Error(msg string) {
	l.fl.Error(msg)
}

// Debug отладочные сообщения.
func (l Logger) Debug(msg, pkg string) {
	l.fl.WithFields(logrus.Fields{
		"package": pkg,
	}).Debug(msg)
}
