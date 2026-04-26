// Package shell provides utilities for executing shell commands with integrated logging support.
//
// The package simplifies running external commands by handling output redirection to a logger,
// supporting verbose mode for debugging, and setting the working directory for command execution.
//
// Key features:
//   - Command execution with customizable working directory
//   - Automatic redirection of stdout/stderr to logger outputs
//   - Verbose mode for detailed execution logging
//   - Support for variable number of command arguments
//
// Typical usage:
//
//	log := logger.NewLogger(os.Stdout, os.Stderr, true)
//	err := shell.Exec("/path/to/workdir", true, "git", log, "status", "--short")
//	if err != nil {
//	    log.Errorf("Command failed: %v", err)
//	}
//
// The Exec function runs the command synchronously and returns any execution error.
// Standard output is only logged when verbose mode is enabled, while standard error
// is always captured through the logger's error writer.
package shell

import (
	"os/exec"

	"github.com/desulaidovich/plea-cli/internal/logger"
)

func Exec(dir string, verbose bool, name string, logger *logger.Logger, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	if verbose {
		logger.Infof("Executing: %s %v in %s", name, args, dir)
		cmd.Stdout = logger.InfoWriter()
	}

	cmd.Stderr = logger.ErrWriter()

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
