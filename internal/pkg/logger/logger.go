package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger defines the interface for logging operations
type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithError(err error) Logger
}

// logrusLogger wraps logrus.Logger to implement our Logger interface
type logrusLogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

// standard logger instance
var stdLogger *logrusLogger

func init() {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{})
	
	// Set log level from environment variable or default to info
	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	log.SetLevel(logLevel)
	
	stdLogger = &logrusLogger{
		logger: log,
		entry:  logrus.NewEntry(log),
	}
}

// GetLogger returns the standard logger instance
func GetLogger() Logger {
	return stdLogger
}

// NewLogger creates a new logger with the specified configuration
func NewLogger(config LoggerConfig) Logger {
	log := logrus.New()
	log.SetOutput(config.Output)
	log.SetFormatter(config.Formatter)
	log.SetLevel(config.Level)
	
	return &logrusLogger{
		logger: log,
		entry:  logrus.NewEntry(log),
	}
}

// LoggerConfig holds the configuration for creating a new logger
type LoggerConfig struct {
	Output    io.Writer
	Formatter logrus.Formatter
	Level     logrus.Level
}

// Info logs at the info level
func (l *logrusLogger) Info(args ...interface{}) {
	if l.entry != nil {
		l.entry.Info(args...)
	} else {
		l.logger.Info(args...)
	}
}

// Infof logs at the info level with formatting
func (l *logrusLogger) Infof(format string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Infof(format, args...)
	} else {
		l.logger.Infof(format, args...)
	}
}

// Warn logs at the warning level
func (l *logrusLogger) Warn(args ...interface{}) {
	if l.entry != nil {
		l.entry.Warn(args...)
	} else {
		l.logger.Warn(args...)
	}
}

// Warnf logs at the warning level with formatting
func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Warnf(format, args...)
	} else {
		l.logger.Warnf(format, args...)
	}
}

// Error logs at the error level
func (l *logrusLogger) Error(args ...interface{}) {
	if l.entry != nil {
		l.entry.Error(args...)
	} else {
		l.logger.Error(args...)
	}
}

// Errorf logs at the error level with formatting
func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Errorf(format, args...)
	} else {
		l.logger.Errorf(format, args...)
	}
}

// Debug logs at the debug level
func (l *logrusLogger) Debug(args ...interface{}) {
	if l.entry != nil {
		l.entry.Debug(args...)
	} else {
		l.logger.Debug(args...)
	}
}

// Debugf logs at the debug level with formatting
func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Debugf(format, args...)
	} else {
		l.logger.Debugf(format, args...)
	}
}

// Fatal logs at the fatal level
func (l *logrusLogger) Fatal(args ...interface{}) {
	if l.entry != nil {
		l.entry.Fatal(args...)
	} else {
		l.logger.Fatal(args...)
	}
}

// Fatalf logs at the fatal level with formatting
func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Fatalf(format, args...)
	} else {
		l.logger.Fatalf(format, args...)
	}
}

// WithField returns a new Logger with the field added
func (l *logrusLogger) WithField(key string, value interface{}) Logger {
	return &logrusLogger{
		logger: l.logger,
		entry:  l.entry.WithField(key, value),
	}
}

// WithFields returns a new Logger with the fields added
func (l *logrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &logrusLogger{
		logger: l.logger,
		entry:  l.entry.WithFields(logrus.Fields(fields)),
	}
}

// WithError returns a new Logger with the error field added
func (l *logrusLogger) WithError(err error) Logger {
	return &logrusLogger{
		logger: l.logger,
		entry:  l.entry.WithError(err),
	}
}