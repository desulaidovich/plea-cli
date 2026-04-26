// Package generator provides functionality for scaffolding new Go service projects.
//
// It creates a standard project structure with sensible defaults including:
//   - Standard Go project layout (cmd/, internal/, pkg/)
//   - Pre-configured logging package with configurable log levels
//   - A runner package for managing application lifecycle
//   - Makefile with common development tasks
//   - Go module initialization with proper module path
//
// The generator uses embedded templates to produce consistent, production-ready
// boilerplate code. It automatically runs 'go mod init' and 'go mod tidy' to
// ensure a valid Go module is created.
//
// Example usage:
//
//	log := logger.New(os.Stdout, true)
//	gen := generator.New(
//	    "my-service",
//	    "github.com/user/my-service",
//	    "./output",
//	    "info",
//	    true,
//	    log,
//	)
//	if err := gen.Do(); err != nil {
//	    log.Fatal(err)
//	}
//
// The generated project includes:
//   - cmd/app/main.go     - Application entry point
//   - internal/app/app.go - Core application logic
//   - pkg/log/log.go      - Structured logging wrapper
//   - pkg/runner/runner.go - Graceful shutdown and signal handling
//   - Makefile            - Build, test, and run targets
package generator

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"plea-cli/internal/logger"
	"plea-cli/internal/shell"
)

//go:embed templates
var templatesFS embed.FS

type TemplateData struct {
	ProjectName string
	Module      string
	Author      string
	LogLevel    string
}

type Generator struct {
	projectName string
	moduleName  string
	outputDir   string
	logLevel    string
	verbose     bool
	logger      *logger.Logger
}

func New(projectName, moduleName, outputDir, logLevel string, verbose bool, logger *logger.Logger) *Generator {
	return &Generator{
		projectName: projectName,
		moduleName:  moduleName,
		outputDir:   outputDir,
		logLevel:    logLevel,
		verbose:     verbose,
		logger:      logger,
	}
}

func (gen *Generator) Do() error {
	serviceDir := filepath.Join(gen.outputDir, gen.projectName)

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
		if gen.verbose {
			gen.logger.Infof("Created directory: %s", dir)
		}
	}

	if err := shell.Exec(serviceDir, gen.verbose, "go", gen.logger, "mod", "init", gen.moduleName); err != nil {
		return fmt.Errorf("go mod init: %w", err)
	}

	entries := []struct{ src, dst string }{
		{"templates/app.go.tmpl", filepath.Join(serviceDir, "internal", "app", "app.go")},
		{"templates/main.go.tmpl", filepath.Join(serviceDir, "cmd", "app", "main.go")},
		{"templates/runner.go.tmpl", filepath.Join(serviceDir, "pkg", "runner", "runner.go")},
		{"templates/log.go.tmpl", filepath.Join(serviceDir, "pkg", "log", "log.go")},
		{"templates/Makefile.tmpl", filepath.Join(serviceDir, "Makefile")},
	}

	data := TemplateData{
		ProjectName: gen.projectName,
		Module:      gen.moduleName,
		LogLevel:    gen.logLevel,
	}

	for _, entrie := range entries {
		if err := gen.render(entrie.src, entrie.dst, data); err != nil {
			return err
		}
		if gen.verbose {
			gen.logger.Infof("Created file: %s", entrie.dst)
		}
	}

	if err := shell.Exec(serviceDir, gen.verbose, "go", gen.logger, "mod", "tidy"); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	return nil
}

func (gen *Generator) render(src, dst string, data any) (err error) {
	content, err := templatesFS.ReadFile(src)
	if err != nil {
		err = fmt.Errorf("failed to read template %q: %w", src, err)
		return
	}

	tmpl, err := template.New(src).Parse(string(content))
	if err != nil {
		err = fmt.Errorf("failed to parse template %q: %w", src, err)
		return
	}

	f, err := os.Create(dst)
	if err != nil {
		err = fmt.Errorf("failed to create %q: %w", dst, err)
		return
	}

	defer func() {
		if errCLose := f.Close(); err != nil {
			err = errors.Join(err, errCLose)
		}
	}()

	if err = tmpl.Execute(f, data); err != nil {
		err = fmt.Errorf("failed to execute template %q: %w", src, err)
		return
	}

	return
}
