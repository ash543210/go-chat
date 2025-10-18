package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type traceKeyType string

const TraceKey traceKeyType = "X-Trace-ID"

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var levelNames = [...]string{
	"DEBUG",
	"INFO ",
	"WARN ",
	"ERROR",
}

type Logger struct {
	mu      sync.Mutex
	level   Level
	logger  *log.Logger
	output  io.Writer
	logFile *os.File
}

// New creates a new logger instance
func New(level Level, toFile bool, filename string) (*Logger, error) {
	var output io.Writer = os.Stdout
	var file *os.File

	if toFile {
		if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
			return nil, err
		}
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		file = f
		output = io.MultiWriter(os.Stdout, f)
	}

	return &Logger{
		level:   level,
		logger:  log.New(output, "", 0),
		output:  output,
		logFile: file,
	}, nil
}

// WithTraceID adds trace ID to the context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceKey, traceID)
}

// extractTraceID fetches trace ID from context if available
func extractTraceID(ctx context.Context) string {
	if v := ctx.Value(TraceKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

// log prints a message with trace ID (if available)
func (l *Logger) log(ctx context.Context, level Level, format string, args ...any) {
	if level < l.level {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	traceID := extractTraceID(ctx)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)

	if traceID != "" {
		l.logger.Printf("[%s] %s [trace_id=%s] - %s", levelNames[level], timestamp, traceID, msg)
	} else {
		l.logger.Printf("[%s] %s - %s", levelNames[level], timestamp, msg)
	}
}

func (l *Logger) Debug(ctx context.Context, format string, args ...any) {
	l.log(ctx, DEBUG, format, args...)
}
func (l *Logger) Info(ctx context.Context, format string, args ...any) {
	l.log(ctx, INFO, format, args...)
}
func (l *Logger) Warn(ctx context.Context, format string, args ...any) {
	l.log(ctx, WARN, format, args...)
}
func (l *Logger) Error(ctx context.Context, format string, args ...any) {
	l.log(ctx, ERROR, format, args...)
}

func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}
