package logging

import (
	"log/slog"
)

// Logger is a generic logging interface "almost" fulfilled by most good logging implementations
//
// Slim wrappers should be created around common logging implementations
// so the actual logger used in an application can be swapped out with ease.
// This interface exists because the "With" and "WithError" methods need to return
// the interface itself, and therefore other libraries will not all implement the same interface,
// but a trivial wrapper can be created for most of them with ease.
type Logger interface {
	With(key string, value any) Logger
	WithError(err error) Logger
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

// slog.Logger implementation of the Logger interface
type SlogLogger struct {
	base *slog.Logger
}

// NewSlogLogger creates a SlogLogger from the provded base slog.Logger
func NewSlogLogger(slogLogger *slog.Logger) *SlogLogger {
	return &SlogLogger{base: slogLogger}
}

// With fulfils Logger.With
func (l *SlogLogger) With(key string, value any) Logger {
	return &SlogLogger{base: l.base.With(key, value)}
}

// WithError fulfils Logger.WithError
func (l *SlogLogger) WithError(err error) Logger {
	return &SlogLogger{base: l.base.With("error", err)}
}

// Debug fulfils Logger.Debug
func (l *SlogLogger) Debug(msg string) {
	l.base.Debug(msg)
}

// Info fulfils Logger.Info
func (l *SlogLogger) Info(msg string) {
	l.base.Info(msg)
}

// Warn fulfils Logger.Warn
func (l *SlogLogger) Warn(msg string) {
	l.base.Warn(msg)
}

// Error fulfils Logger.Error
func (l *SlogLogger) Error(msg string) {
	l.base.Error(msg)
}
