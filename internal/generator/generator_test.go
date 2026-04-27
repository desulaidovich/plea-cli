package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/desulaidovich/plea-cli/internal/logger"
	"github.com/desulaidovich/plea-cli/internal/manifest"
)

func TestGenerator_Do(t *testing.T) {
	t.Run("creates project structure", func(t *testing.T) {
		tempDir := t.TempDir()
		cfg := manifest.Config{
			Name:     "test-project",
			Module:   "github.com/test/test-project",
			Output:   tempDir,
			LogLevel: "debug",
			Verbose:  false,
		}

		err := New(cfg, logger.New(&bytes.Buffer{}, logger.Debug, 0)).Do()
		require.NoError(t, err)

		expectedDirs := []string{
			filepath.Join(tempDir, cfg.Name),
			filepath.Join(tempDir, cfg.Name, "cmd", "app"),
			filepath.Join(tempDir, cfg.Name, "internal", "app"),
			filepath.Join(tempDir, cfg.Name, "pkg", "log"),
			filepath.Join(tempDir, cfg.Name, "pkg", "runner"),
		}
		for _, dir := range expectedDirs {
			_, err = os.Stat(dir)
			assert.NoError(t, err, "directory should exist: %s", dir)
		}

		expectedFiles := []string{
			"internal/app/app.go",
			"cmd/app/main.go",
			"pkg/runner/runner.go",
			"pkg/log/log.go",
			"Makefile",
			"go.mod",
		}
		for _, file := range expectedFiles {
			_, err = os.Stat(filepath.Join(tempDir, cfg.Name, file))
			assert.NoError(t, err, "file should exist: %s", file)
		}
	})

	t.Run("fails when project already exists", func(t *testing.T) {
		tempDir := t.TempDir()
		cfg := manifest.Config{
			Name:   "test-project",
			Module: "github.com/test/test-project",
			Output: tempDir,
		}
		gen := New(cfg, logger.New(&bytes.Buffer{}, logger.Debug, 0))

		require.NoError(t, gen.Do())
		err := gen.Do()
		assert.Error(t, err)
	})
}
