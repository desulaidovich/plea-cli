// Package generator scaffolds new Go service projects from embedded templates.
//
// Example usage:
//
//	gen := generator.New(manifest.Config{
//	    Name:     "my-service",
//	    Module:   "github.com/user/my-service",
//	    Output:   "./output",
//	    LogLevel: "info",
//	    Verbose:  true,
//	}, logger.New(os.Stdout, logger.Debug, log.Ltime))
//	if err := gen.Do(); err != nil {
//	    log.Fatal(err)
//	}
//
// The generated project includes:
//   - cmd/app/main.go      - Application entry point
//   - internal/app/app.go  - Core application logic
//   - pkg/log/log.go       - Structured logging wrapper
//   - pkg/runner/runner.go - Graceful shutdown and signal handling
//   - Makefile             - Build, test, and run targets
package generator

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/desulaidovich/plea-cli/internal/logger"
	"github.com/desulaidovich/plea-cli/internal/manifest"
	"github.com/desulaidovich/plea-cli/internal/shell"
)

//go:embed templates
var templatesFS embed.FS

type templateData struct {
	ProjectName string
	Module      string
	LogLevel    string
}

type Generator struct {
	cfg    manifest.Config
	logger *logger.Logger
}

func New(cfg manifest.Config, logger *logger.Logger) *Generator {
	return &Generator{cfg: cfg, logger: logger}
}

func (gen *Generator) Do() error {
	serviceDir := filepath.Join(gen.cfg.Output, gen.cfg.Name)

	dirs := []string{
		serviceDir,
		filepath.Join(serviceDir, "cmd", "app"),
		filepath.Join(serviceDir, "internal", "app"),
		filepath.Join(serviceDir, "pkg", "log"),
		filepath.Join(serviceDir, "pkg", "runner"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create %q: %w", dir, err)
		}
		if gen.cfg.Verbose {
			gen.logger.Infof("Created directory: %s", dir)
		}
	}

	if err := shell.Exec(serviceDir, gen.cfg.Verbose, "go", gen.logger, "mod", "init", gen.cfg.Module); err != nil {
		return fmt.Errorf("go mod init: %w", err)
	}

	entries := []struct{ src, dst string }{
		{"templates/app.go.tmpl", filepath.Join(serviceDir, "internal", "app", "app.go")},
		{"templates/main.go.tmpl", filepath.Join(serviceDir, "cmd", "app", "main.go")},
		{"templates/runner.go.tmpl", filepath.Join(serviceDir, "pkg", "runner", "runner.go")},
		{"templates/log.go.tmpl", filepath.Join(serviceDir, "pkg", "log", "log.go")},
		{"templates/Makefile.tmpl", filepath.Join(serviceDir, "Makefile")},
	}

	data := templateData{
		ProjectName: gen.cfg.Name,
		Module:      gen.cfg.Module,
		LogLevel:    gen.cfg.LogLevel,
	}

	for _, entry := range entries {
		if err := gen.render(entry.src, entry.dst, data); err != nil {
			return err
		}
		if gen.cfg.Verbose {
			gen.logger.Infof("Created file: %s", entry.dst)
		}
	}

	if err := shell.Exec(serviceDir, gen.cfg.Verbose, "go", gen.logger, "mod", "tidy"); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	return nil
}

func (gen *Generator) render(src, dst string, data any) error {
	content, err := templatesFS.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read template %q: %w", src, err)
	}

	tmpl, err := template.New(src).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %q: %w", src, err)
	}

	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create %q: %w", dst, err)
	}

	execErr := tmpl.Execute(f, data)
	closeErr := f.Close()

	if execErr != nil {
		return fmt.Errorf("failed to execute template %q: %w", src, execErr)
	}
	if closeErr != nil {
		return fmt.Errorf("failed to close %q: %w", dst, closeErr)
	}

	return nil
}
