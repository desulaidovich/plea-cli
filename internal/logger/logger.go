package logger

import (
	"io"
	"log"
	"os"
)

var Default *Logger

func init() {
	Default = New(os.Stdout)
}

type Logger struct {
	debugLog *log.Logger
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
	fatalLog *log.Logger
}

func New(output io.Writer) *Logger {
	if output == nil {
		output = os.Stdout
	}

	return &Logger{
		debugLog: log.New(output, "🔍 ", log.Ltime),
		infoLog:  log.New(output, "ℹ️  ", log.Ltime),
		warnLog:  log.New(output, "⚠️  ", log.Ltime),
		errorLog: log.New(output, "❌ ", log.Ltime),
		fatalLog: log.New(output, "💀 ", log.Ltime),
	}
}

func (l *Logger) InfoWriter() io.Writer {
	return l.infoLog.Writer()
}

func (l *Logger) ErrWriter() io.Writer {
	return l.errorLog.Writer()
}

func (l *Logger) Debug(v ...any) {
	l.debugLog.Print(v...)
}

func (l *Logger) Debugf(format string, v ...any) {
	l.debugLog.Printf(format, v...)
}

func (l *Logger) Info(v ...any) {
	l.infoLog.Print(v...)
}

func (l *Logger) Infof(format string, v ...any) {
	l.infoLog.Printf(format, v...)
}

func (l *Logger) Warn(v ...any) {
	l.warnLog.Print(v...)
}

func (l *Logger) Warnf(format string, v ...any) {
	l.warnLog.Printf(format, v...)
}

func (l *Logger) Error(v ...any) {
	l.errorLog.Print(v...)
}

func (l *Logger) Errorf(format string, v ...any) {
	l.errorLog.Printf(format, v...)
}

func (l *Logger) Fatal(v ...any) {
	l.fatalLog.Fatal(v...)
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.fatalLog.Fatalf(format, v...)
}
