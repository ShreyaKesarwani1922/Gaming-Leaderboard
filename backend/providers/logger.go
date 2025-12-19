package providers

import (
	"fmt"
	"log"
	"os"
)

// LoggerInterface defines the logging interface for the core package
type LoggerInterface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// ConsoleLogger is a simple console logger that implements LoggerInterface
type ConsoleLogger struct {
	logger *log.Logger
}

// NewConsoleLogger creates a new ConsoleLogger instance
func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

func (l *ConsoleLogger) Debug(args ...interface{}) {
	l.logger.SetPrefix("[DEBUG] ")
	l.logger.Output(2, fmt.Sprint(args...))
}

func (l *ConsoleLogger) Info(args ...interface{}) {
	l.logger.SetPrefix("[INFO] ")
	l.logger.Output(2, fmt.Sprint(args...))
}

func (l *ConsoleLogger) Warn(args ...interface{}) {
	l.logger.SetPrefix("[WARN] ")
	l.logger.Output(2, fmt.Sprint(args...))
}

func (l *ConsoleLogger) Error(args ...interface{}) {
	l.logger.SetPrefix("[ERROR] ")
	l.logger.Output(2, fmt.Sprint(args...))
}

func (l *ConsoleLogger) Fatal(args ...interface{}) {
	l.logger.SetPrefix("[FATAL] ")
	l.logger.Output(2, fmt.Sprint(args...))
	os.Exit(1)
}

func (l *ConsoleLogger) Debugf(format string, args ...interface{}) {
	l.logger.SetPrefix("[DEBUG] ")
	l.logger.Output(2, fmt.Sprintf(format, args...))
}

func (l *ConsoleLogger) Infof(format string, args ...interface{}) {
	l.logger.SetPrefix("[INFO] ")
	l.logger.Output(2, fmt.Sprintf(format, args...))
}

func (l *ConsoleLogger) Warnf(format string, args ...interface{}) {
	l.logger.SetPrefix("[WARN] ")
	l.logger.Output(2, fmt.Sprintf(format, args...))
}

func (l *ConsoleLogger) Errorf(format string, args ...interface{}) {
	l.logger.SetPrefix("[ERROR] ")
	l.logger.Output(2, fmt.Sprintf(format, args...))
}

func (l *ConsoleLogger) Fatalf(format string, args ...interface{}) {
	l.logger.SetPrefix("[FATAL] ")
	l.logger.Output(2, fmt.Sprintf(format, args...))
	os.Exit(1)
}
