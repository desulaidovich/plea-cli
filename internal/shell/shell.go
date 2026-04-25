package shell

import (
	"os/exec"
	"plea-cli/internal/logger"
)

func Exec(dir string, verbose bool, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	if verbose {
		logger.Default.Infof("Executing: %s %v in %s", name, args, dir)
		cmd.Stdout = logger.Default.InfoWriter()
		cmd.Stderr = logger.Default.ErrWriter()
	}

	if err := cmd.Run(); err != nil {
		logger.Default.Errorf("Executing: %s %v in %s: %v", name, args, dir, err)
		return err
	}
	return nil
}
