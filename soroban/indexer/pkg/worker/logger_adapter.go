package worker

import "github.com/sirupsen/logrus"

// LogrusAdapter adapts logrus.Entry to the Logger interface
type LogrusAdapter struct {
	entry *logrus.Entry
}

// NewLogrusAdapter creates a new adapter from a logrus logger
func NewLogrusAdapter(logger *logrus.Logger) *LogrusAdapter {
	return &LogrusAdapter{
		entry: logrus.NewEntry(logger),
	}
}

// WithError adds an error to the logger context
func (l *LogrusAdapter) WithError(err error) Logger {
	return &LogrusAdapter{
		entry: l.entry.WithError(err),
	}
}

// WithField adds a field to the logger context
func (l *LogrusAdapter) WithField(key string, value interface{}) Logger {
	return &LogrusAdapter{
		entry: l.entry.WithField(key, value),
	}
}

// WithFields adds multiple fields to the logger context
func (l *LogrusAdapter) WithFields(fields map[string]interface{}) Logger {
	return &LogrusAdapter{
		entry: l.entry.WithFields(fields),
	}
}

// Warn logs a warning message
func (l *LogrusAdapter) Warn(msg string) {
	l.entry.Warn(msg)
}

// Error logs an error message
func (l *LogrusAdapter) Error(msg string) {
	l.entry.Error(msg)
}

// Info logs an info message
func (l *LogrusAdapter) Info(msg string) {
	l.entry.Info(msg)
}
