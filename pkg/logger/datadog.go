package logger

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// DatadogLogger extends logrus.Logger with Datadog integration
type DatadogLogger struct {
	*logrus.Logger
}

// NewDatadogLogger creates a new Datadog-integrated logger
func NewDatadogLogger() *DatadogLogger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	return &DatadogLogger{Logger: logger}
}

// WithContext creates a logger with Datadog trace context
func (l *DatadogLogger) WithContext(ctx context.Context) *logrus.Entry {
	entry := l.Logger.WithContext(ctx)

	// Add Datadog trace information if available
	if span, ok := tracer.SpanFromContext(ctx); ok {
		entry = entry.WithFields(logrus.Fields{
			"dd.trace_id": span.Context().TraceID(),
			"dd.span_id":  span.Context().SpanID(),
		})
	}

	return entry
}

// WithSpan creates a logger with a specific span
func (l *DatadogLogger) WithSpan(span tracer.Span) *logrus.Entry {
	entry := l.Logger.WithFields(logrus.Fields{
		"dd.trace_id": span.Context().TraceID(),
		"dd.span_id":  span.Context().SpanID(),
	})

	return entry
}

// LogWithSpan logs a message with span information
func (l *DatadogLogger) LogWithSpan(span tracer.Span, level logrus.Level, msg string, fields logrus.Fields) {
	entry := l.WithSpan(span)
	if fields != nil {
		entry = entry.WithFields(fields)
	}
	entry.Log(level, msg)
}

// LogWithContext logs a message with context information
func (l *DatadogLogger) LogWithContext(ctx context.Context, level logrus.Level, msg string, fields logrus.Fields) {
	entry := l.WithContext(ctx)
	if fields != nil {
		entry = entry.WithFields(fields)
	}
	entry.Log(level, msg)
}

// InfoWithSpan logs an info message with span
func (l *DatadogLogger) InfoWithSpan(span tracer.Span, msg string, fields logrus.Fields) {
	l.LogWithSpan(span, logrus.InfoLevel, msg, fields)
}

// ErrorWithSpan logs an error message with span
func (l *DatadogLogger) ErrorWithSpan(span tracer.Span, msg string, fields logrus.Fields) {
	l.LogWithSpan(span, logrus.ErrorLevel, msg, fields)
}

// InfoWithContext logs an info message with context
func (l *DatadogLogger) InfoWithContext(ctx context.Context, msg string, fields logrus.Fields) {
	l.LogWithContext(ctx, logrus.InfoLevel, msg, fields)
}

// ErrorWithContext logs an error message with context
func (l *DatadogLogger) ErrorWithContext(ctx context.Context, msg string, fields logrus.Fields) {
	l.LogWithContext(ctx, logrus.ErrorLevel, msg, fields)
}

// Info logs an info message
func (l *DatadogLogger) Info(msg string, fields logrus.Fields) {
	if fields != nil {
		l.Logger.WithFields(fields).Info(msg)
	} else {
		l.Logger.Info(msg)
	}
}

// Error logs an error message
func (l *DatadogLogger) Error(msg string, fields logrus.Fields) {
	if fields != nil {
		l.Logger.WithFields(fields).Error(msg)
	} else {
		l.Logger.Error(msg)
	}
}

// SetLevel sets the log level
func (l *DatadogLogger) SetLevel(level logrus.Level) {
	l.Logger.SetLevel(level)
}
