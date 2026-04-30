package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerLevels(t *testing.T) {
	tests := []struct {
		name       string
		level      Level
		setUp      func(*Logger)
		wantOutput bool
	}{
		{
			name:  "when level=DEBUG",
			level: Debug,
			setUp: func(l *Logger) {
				l.Debug("test")
			},
			wantOutput: true,
		},
		{
			name:  "when level=DEBUG (info)",
			level: Debug,
			setUp: func(l *Logger) {
				l.Info("test")
			},
			wantOutput: true,
		},
		{
			name:  "when level=INFO",
			level: Info,
			setUp: func(l *Logger) {
				l.Info("test")
			},
			wantOutput: true,
		},
		{
			name:  "when level=INFO (debug)",
			level: Info,
			setUp: func(l *Logger) {
				l.Debug("test")
			},
			wantOutput: false,
		},
		{
			name:  "when level=WARN",
			level: Warn,
			setUp: func(l *Logger) {
				l.Warn("test")
			},
			wantOutput: true,
		},
		{
			name:  "when level=WARN (info)",
			level: Warn,
			setUp: func(l *Logger) {
				l.Info("test")
			},
			wantOutput: false,
		},
		{
			name:  "when level=ERROR",
			level: Error,
			setUp: func(l *Logger) {
				l.Error("test")
			},
			wantOutput: true,
		},
		{
			name:  "when level=ERROR (warn)",
			level: Error,
			setUp: func(l *Logger) {
				l.Warn("test")
			},
			wantOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buff bytes.Buffer
			logger := New(&buff, tt.level, 0)

			tt.setUp(logger)

			output := buff.String()
			if tt.wantOutput {
				assert.NotEmpty(t, output)
			} else {
				assert.Empty(t, output)
			}
		})
	}
}
