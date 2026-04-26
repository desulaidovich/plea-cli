package shell

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"plea-cli/internal/logger"
	"testing"
)

func TestExec(t *testing.T) {
	tests := []struct {
		name      string
		command   string
		args      []string
		verbose   bool
		dir       string
		setupFunc func(t *testing.T) string
		wantErr   bool
		errCheck  func(error) bool
	}{
		{
			name:    "successful command with verbose",
			command: "echo",
			args:    []string{"test message"},
			verbose: true,
			dir:     ".",
			wantErr: false,
		},
		{
			name:    "successful command without verbose",
			command: "echo",
			args:    []string{"silent message"},
			verbose: false,
			dir:     ".",
			wantErr: false,
		},
		{
			name:    "command with error (non-zero exit code)",
			command: "sh",
			args:    []string{"-c", "exit 1"},
			verbose: true,
			dir:     ".",
			wantErr: true,
		},
		{
			name:    "command not found",
			command: "nonexistent_command_xyz_123",
			args:    []string{},
			verbose: false,
			dir:     ".",
			wantErr: true,
			errCheck: func(err error) bool {
				var execErr *exec.Error
				return errors.As(err, &execErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuf bytes.Buffer
			log := logger.New(&logBuf, logger.Debug, log.LstdFlags)

			err := Exec(tt.dir, tt.verbose, tt.command, log, tt.args...)

			if (err != nil) != tt.wantErr {
				t.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.errCheck != nil && err != nil {
				if !tt.errCheck(err) {
					t.Errorf("Expected specific error type, got %v", err)
				}
			}
		})
	}
}
