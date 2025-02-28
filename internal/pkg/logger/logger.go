package logger

import (
    "github.com/sirupsen/logrus"
    "os"
)

var log = logrus.New()

func init() {
    log.SetOutput(os.Stdout)
    log.SetFormatter(&logrus.JSONFormatter{})
    
    // Set log level from environment variable or default to info
    logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
    if err != nil {
        logLevel = logrus.InfoLevel
    }
    log.SetLevel(logLevel)
}

// GetLogger returns the logger instance
func GetLogger() *logrus.Logger {
    return log
}