// Package logger provides structured logging functionality.
package logger

import (
	"os"
	"strings"
	"time"

	"typesense-mcp-server/internal/constants"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus.Logger with additional functionality
type Logger struct {
	*logrus.Logger
}

// Config holds the logger configuration
type Config struct {
	Level  string
	Format string
	Output string
}

// New creates a new Logger instance
func New(config Config) *Logger {
	logger := logrus.New()
	
	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)
	
	// Set log format
	switch strings.ToLower(config.Format) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "function",
			},
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		})
	}
	
	// Set output
	switch strings.ToLower(config.Output) {
	case "file":
		file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.SetOutput(file)
		}
	default:
		logger.SetOutput(os.Stdout)
	}
	
	return &Logger{Logger: logger}
}

// NewDefault creates a logger with default configuration
func NewDefault() *Logger {
	config := Config{
		Level:  getEnvOrDefault(constants.EnvLogLevel, "info"),
		Format: getEnvOrDefault(constants.EnvLogFormat, "text"),
		Output: "stdout",
	}
	return New(config)
}

// WithComponent adds a component field to the logger
func (l *Logger) WithComponent(component string) *logrus.Entry {
	return l.WithField(constants.LogFieldComponent, component)
}

// WithOperation adds an operation field to the logger
func (l *Logger) WithOperation(operation string) *logrus.Entry {
	return l.WithField(constants.LogFieldOperation, operation)
}

// WithRequestID adds a request ID field to the logger
func (l *Logger) WithRequestID(requestID string) *logrus.Entry {
	return l.WithField(constants.LogFieldRequestID, requestID)
}

// WithCollection adds a collection field to the logger
func (l *Logger) WithCollection(collection string) *logrus.Entry {
	return l.WithField(constants.LogFieldCollection, collection)
}

// WithError adds an error field to the logger
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.WithField(constants.LogFieldError, err)
}

// WithDuration adds a duration field to the logger
func (l *Logger) WithDuration(duration time.Duration) *logrus.Entry {
	return l.WithField(constants.LogFieldDuration, duration.Milliseconds())
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *logrus.Entry {
	return l.Logger.WithFields(logrus.Fields(fields))
}

// TimedOperation logs the start and end of an operation with duration
func (l *Logger) TimedOperation(component, operation string, fn func() error) error {
	start := time.Now()
	entry := l.WithComponent(component).WithField(constants.LogFieldOperation, operation)
	
	entry.Info("operation started")
	
	err := fn()
	duration := time.Since(start)
	
	if err != nil {
		entry.WithError(err).WithField(constants.LogFieldDuration, duration.Milliseconds()).Error("operation failed")
	} else {
		entry.WithField(constants.LogFieldDuration, duration.Milliseconds()).Info("operation completed")
	}
	
	return err
}

// LogSearchOperation logs a search operation with relevant context
func (l *Logger) LogSearchOperation(collection, query string, page, perPage, resultCount int, duration time.Duration, err error) {
	entry := l.WithComponent(constants.ComponentHandler).
		WithField(constants.LogFieldOperation, "search").
		WithField(constants.LogFieldCollection, collection).
		WithField(constants.LogFieldQuery, query).
		WithField(constants.LogFieldPage, page).
		WithField(constants.LogFieldPageSize, perPage).
		WithField(constants.LogFieldResultCount, resultCount).
		WithField(constants.LogFieldDuration, duration.Milliseconds())
	
	if err != nil {
		entry.WithError(err).Error("search operation failed")
	} else {
		entry.Info("search operation completed")
	}
}

// LogCollectionOperation logs a collection operation with relevant context
func (l *Logger) LogCollectionOperation(operation string, count int, duration time.Duration, err error) {
	entry := l.WithComponent(constants.ComponentHandler).
		WithField(constants.LogFieldOperation, operation).
		WithField("collection_count", count).
		WithField(constants.LogFieldDuration, duration.Milliseconds())
	
	if err != nil {
		entry.WithError(err).Error("collection operation failed")
	} else {
		entry.Info("collection operation completed")
	}
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
