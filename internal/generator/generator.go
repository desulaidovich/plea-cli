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
}

func New(projectName, moduleName, outputDir, logLevel string, verbose bool) *Generator {
	return &Generator{
		projectName: projectName,
		moduleName:  moduleName,
		outputDir:   outputDir,
		logLevel:    logLevel,
		verbose:     verbose,
	}
}

func (gen *Generator) Do() error {
	serviceDir := gen.outputDir + "/generated/" + gen.projectName

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
			logger.Default.Infof("Created directory: %s", dir)
		}
	}

	if err := shell.Exec(serviceDir, gen.verbose, "go", "mod", "init", gen.moduleName); err != nil {
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
			logger.Default.Infof("Created file: %s", entrie.dst)
		}
	}

	if err := shell.Exec(serviceDir, gen.verbose, "go", "mod", "tidy"); err != nil {
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
