package logger

import (
	"log"
	"os"
)

type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func New() *Logger {
	return &Logger{
		infoLog:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		errorLog: log.New(os.Stdout, "[ERROR] ", log.LstdFlags),
	}
}

func (l *Logger) Info(format string, args ...any) {
	l.infoLog.Printf(format, args...)
}

func (l *Logger) Error(format string, args ...any) {
	l.errorLog.Printf(format, args...)
}
