package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"plea-cli/internal/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerator_Do(t *testing.T) {
	type args struct {
		projectName string
		moduleName  string
		outputDir   string
		logLevel    string
		verbose     bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "temporary directory",
			args: args{
				projectName: "test-project",
				moduleName:  "github.com/test/test-project",
				outputDir:   "test",
				logLevel:    "debug",
				verbose:     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", tt.args.outputDir)
			assert.NoError(t, err)
			defer func() {
				err = os.RemoveAll(tempDir)
				assert.NoError(t, err)
			}()

			var buff bytes.Buffer

			gen := New(
				tt.args.projectName,
				tt.args.moduleName,
				tempDir,
				tt.args.logLevel,
				tt.args.verbose,
				logger.New(&buff, logger.Debug, 0),
			)

			err = gen.Do()
			assert.NoError(t, err)

			expectedDirs := []string{
				filepath.Join(tempDir, tt.args.projectName),
				filepath.Join(tempDir, tt.args.projectName, "cmd", "app"),
				filepath.Join(tempDir, tt.args.projectName, "internal", "app"),
				filepath.Join(tempDir, tt.args.projectName, "pkg", "log"),
				filepath.Join(tempDir, tt.args.projectName, "pkg", "runner"),
			}

			for _, dir := range expectedDirs {
				_, err = os.Stat(dir)
				assert.NoError(t, err, "Directory should exist: %s", dir)
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
				filePath := filepath.Join(tempDir, tt.args.projectName, file)
				_, err = os.Stat(filePath)
				assert.NoError(t, err, "File should exist: %s", file)
			}
		})
	}
}
func TestGenerator_render(t *testing.T) {
	type args struct {
		src  string
		dst  string
		data any
	}
	tests := []struct {
		name    string
		gen     *Generator
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.gen.render(tt.args.src, tt.args.dst, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Generator.render() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
