package helpers

import (
	"log"
	"os"
	"sync"
)

type Logger struct {
	logger *log.Logger
	mu     sync.Mutex
}

func NewLogger() (*Logger, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	return &Logger{
		logger: logger,
	}, nil
}

func (l *Logger) LogMessage(message string) {
	l.mu.Lock()
	l.logger.Printf("INFO: %s", message)
	l.mu.Unlock()
}

func (l *Logger) LogInfof(format string, v ...interface{}) {
	l.mu.Lock()
	l.logger.Printf("INFO: "+format, v...)
	l.mu.Unlock()
}

func (l *Logger) LogErrorMessage(message string) {
	l.mu.Lock()
	l.logger.Printf("ERROR: %s", message)
	l.mu.Unlock()
}

func (l *Logger) LogErrorf(format string, v ...interface{}) {
	l.mu.Lock()
	l.logger.Printf("ERROR: "+format, v...)
	l.mu.Unlock()
}

func (l *Logger) LogWarnMessage(message string) {
	l.mu.Lock()
	l.logger.Printf("WARNING: %s", message)
	l.mu.Unlock()
}

func (l *Logger) LogWarnf(format string, v ...interface{}) {
	l.mu.Lock()
	l.logger.Printf("WARNING: "+format, v...)
	l.mu.Unlock()
}

func (l *Logger) LogDebugMessage(message string) {
	l.mu.Lock()
	l.logger.Printf("DEBUG: %s", message)
	l.mu.Unlock()
}

func (l *Logger) LogDebugf(format string, v ...interface{}) {
	l.mu.Lock()
	l.logger.Printf("DEBUG: "+format, v...)
	l.mu.Unlock()
}