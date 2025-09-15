package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
)

// Logger wraps zap.Logger to provide a consistent logging interface
type Logger struct {
	*zap.Logger
}

// Custom is the global logger instance that can be used throughout the application
var Custom *Logger

// New creates a new Logger instance using environment-based configuration
func New() *Logger {
	env := os.Getenv("ENVIRONMENT")
	if env == "development" || env == "dev" {
		return NewDevelopment()
	}
	return NewProduction()
}

// NewProduction creates a new Logger instance optimized for production use.
// Uses JSON encoding and focuses on performance and structured output.
func NewProduction() *Logger {
	logger, err := zap.NewProduction(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatalf("can't initialize zap production logger: %v", err)
	}
	return &Logger{logger}
}

// NewDevelopment creates a new Logger instance optimized for development use.
// Uses console encoding with colors and human-readable format.
func NewDevelopment() *Logger {
	logger, err := zap.NewDevelopment(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatalf("can't initialize zap development logger: %v", err)
	}
	return &Logger{logger}
}

// Info logs a message at Info level with structured key-value pairs
func (l *Logger) Info(msg string, fields ...any) {
	l.Logger.Sugar().Infow(msg, fields...)
}

// Error logs a message at Error level with structured key-value pairs
func (l *Logger) Error(msg string, fields ...any) {
	l.Logger.Sugar().Errorw(msg, fields...)
}

// Warn logs a message at Warn level with structured key-value pairs
func (l *Logger) Warn(msg string, fields ...any) {
	l.Logger.Sugar().Warnw(msg, fields...)
}

// Debug logs a message at Debug level with structured key-value pairs
func (l *Logger) Debug(msg string, fields ...any) {
	l.Logger.Sugar().Debugw(msg, fields...)
}

// Fatal logs a message at Fatal level with structured key-value pairs and exits the program
func (l *Logger) Fatal(msg string, fields ...any) {
	l.Logger.Sugar().Fatalw(msg, fields...)
}

// Sync flushes any buffered log entries.
// Should be called before the application exits to ensure all logs are written.
func (l *Logger) Sync() {
	if err := l.Logger.Sync(); err != nil {
		log.Printf("failed to sync logger: %v", err)
	}
}

// Close performs cleanup operations on the logger.
// It syncs any remaining logs and can be used in defer statements for proper resource management.
func (l *Logger) Close() error {
	l.Sync()
	return nil
}

// InitGlobal initializes the global Custom logger instance
func InitGlobal() {
	Custom = New()
}

// InitGlobalProduction initializes the global Custom logger instance for production
func InitGlobalProduction() {
	Custom = NewProduction()
}

// InitGlobalDevelopment initializes the global Custom logger instance for development
func InitGlobalDevelopment() {
	Custom = NewDevelopment()
}

// Package-level logging functions that use the global Custom logger

// Info logs a message at Info level with structured key-value pairs
func Info(msg string, fields ...any) {
	Custom.Info(msg, fields...)
}

// Error logs a message at Error level with structured key-value pairs
func Error(msg string, fields ...any) {
	Custom.Error(msg, fields...)
}

// Warn logs a message at Warn level with structured key-value pairs
func Warn(msg string, fields ...any) {
	Custom.Warn(msg, fields...)
}

// Debug logs a message at Debug level with structured key-value pairs
func Debug(msg string, fields ...any) {
	Custom.Debug(msg, fields...)
}

// Fatal logs a message at Fatal level with structured key-value pairs and exits the program
func Fatal(msg string, fields ...any) {
	Custom.Fatal(msg, fields...)
}