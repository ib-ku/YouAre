package logger

import (
	"log"
	"os"
	"sync"
)

type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	warnLog  *log.Logger
	debugLog *log.Logger
	mu       sync.Mutex
}

var (
	instance *Logger
	once     sync.Once
)

func GetLogger() *Logger {
	once.Do(func() {
		flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
		instance = &Logger{
			infoLog: log.New(os.Stdout,
				"INFO: ",
				flags),
			errorLog: log.New(os.Stderr,
				"ERROR: ",
				flags),
			warnLog: log.New(os.Stdout,
				"WARN: ",
				flags),
			debugLog: log.New(os.Stdout,
				"DEBUG: ",
				flags),
		}
	})
	return instance
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.infoLog.Printf(format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.errorLog.Printf(format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.warnLog.Printf(format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.debugLog.Printf(format, v...)
}
