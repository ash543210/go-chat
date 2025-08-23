package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

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
		var err error
		file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		output = io.MultiWriter(os.Stdout, file)
	}

	return &Logger{
		level:   level,
		logger:  log.New(output, "", 0),
		output:  output,
		logFile: file,
	}, nil
}

// log prints a message if the level is appropriate
func (l *Logger) log(level Level, format string, args ...any) {
	if level < l.level {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	l.logger.Printf("[%s] %s - %s", levelNames[level], timestamp, msg)
}

func (l *Logger) Debug(format string, args ...any) { l.log(DEBUG, format, args...) }
func (l *Logger) Info(format string, args ...any)  { l.log(INFO, format, args...) }
func (l *Logger) Warn(format string, args ...any)  { l.log(WARN, format, args...) }
func (l *Logger) Error(format string, args ...any) { l.log(ERROR, format, args...) }

func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}
