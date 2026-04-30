// Package logger provides a leveled logging wrapper around the standard log package.
//
// It supports five log levels: Debug, Info, Warn, Error, and Fatal, allowing
// granular control over log output based on severity. Each level has a
// corresponding colored or emoji-based prefix (🔍, ℹ️, ⚠️, ❌, 💀) for visual
// distinction in console output.
//
// The logger can be configured with:
//   - An output writer (defaults to os.Stdout if nil)
//   - A minimum log level (messages below this level are suppressed)
//   - Standard log flags (date, time, file location, etc.)
//
// Two convenience methods provide access to underlying io.Writer instances:
//   - InfoWriter() returns the writer for info-level logs
//   - ErrWriter() returns the writer for error-level logs
//
// Example usage:
//
//	// Create a logger with INFO level and timestamps
//	log := logger.New(os.Stdout, logger.Info, log.Ldate|log.Ltime)
//
//	log.Debug("This won't appear")  // level is Debug < Info
//	log.Info("Application started")
//	log.Warnf("Disk usage: %d%%", 85)
//	log.Error("Connection failed")
//	// log.Fatal("Critical error")  // exits the program
//
// Note: Fatal and Fatalf call os.Exit(1) after writing the log message,
// terminating the program immediately.
package logger

import (
	"io"
	"log"
	"os"
)

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

type Logger struct {
	level Level
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
	fatal *log.Logger
}

func New(output io.Writer, level Level, flag int) *Logger {
	if output == nil {
		output = os.Stdout
	}

	return &Logger{
		level: level,
		debug: log.New(output, "🔍 ", flag),
		info:  log.New(output, "ℹ️  ", flag),
		warn:  log.New(output, "⚠️  ", flag),
		error: log.New(output, "❌ ", flag),
		fatal: log.New(output, "💀 ", flag),
	}
}

func (l *Logger) InfoWriter() io.Writer {
	return l.info.Writer()
}

func (l *Logger) ErrWriter() io.Writer {
	return l.error.Writer()
}

func (l *Logger) Debug(v ...any) {
	if l.level <= Debug {
		l.debug.Print(v...)
	}
}

func (l *Logger) Info(v ...any) {
	if l.level <= Info {
		l.info.Print(v...)
	}
}

func (l *Logger) Warn(v ...any) {
	if l.level <= Warn {
		l.warn.Print(v...)
	}
}

func (l *Logger) Error(v ...any) {
	if l.level <= Error {
		l.error.Print(v...)
	}
}

func (l *Logger) Fatal(v ...any) {
	l.fatal.Fatal(v...)
}

func (l *Logger) Debugf(format string, v ...any) {
	if l.level <= Debug {
		l.debug.Printf(format, v...)
	}
}

func (l *Logger) Infof(format string, v ...any) {
	if l.level <= Info {
		l.info.Printf(format, v...)
	}
}

func (l *Logger) Warnf(format string, v ...any) {
	if l.level <= Warn {
		l.warn.Printf(format, v...)
	}
}

func (l *Logger) Errorf(format string, v ...any) {
	if l.level <= Error {
		l.error.Printf(format, v...)
	}
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.fatal.Fatalf(format, v...)
}
