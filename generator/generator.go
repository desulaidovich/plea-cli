// Package generator scaffolds new Go projects from a blueprint and embedded templates.
// It reads the project structure from blueprint.yaml and renders each file through a
// matching text/template.
package generator

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/desulaidovich/plea-cli/blueprint"
)

//go:embed templates
var templatesFS embed.FS

type Logger interface {
	Debug(v ...any)
	Debugf(format string, v ...any)
	ErrWriter() io.Writer
	Error(v ...any)
	Errorf(format string, v ...any)
	Fatal(v ...any)
	Fatalf(format string, v ...any)
	Info(v ...any)
	InfoWriter() io.Writer
	Infof(format string, v ...any)
	Warn(v ...any)
	Warnf(format string, v ...any)
}

type templateData struct {
	ProjectName string
	ModulePath  string
}

type Generator struct {
	projectName    string
	repositoryPath string
	verbose        bool
	logger         Logger
}

func New(projectName, repositoryPath string, logger Logger, verbose bool) *Generator {
	return &Generator{
		projectName:    projectName,
		repositoryPath: repositoryPath,
		verbose:        verbose,
		logger:         logger,
	}
}

func (gen *Generator) Build() error {
	if err := validateProjectName(gen.projectName); err != nil {
		return err
	}

	projectDir := filepath.Join(".", gen.projectName)

	if _, err := os.Stat(projectDir); err == nil {
		return fmt.Errorf("directory %q already exists", projectDir)
	}

	bp, err := blueprint.Read()
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", blueprint.RootFile, err)
	}

	if strings.Contains(bp.Module, "{{repositoryPath}}") && gen.repositoryPath == "" {
		return fmt.Errorf("blueprint requires --repo: module path contains {{repositoryPath}}")
	}

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project dir %q: %w", projectDir, err)
	}

	data := templateData{
		ProjectName: gen.projectName,
		ModulePath:  bp.ModulePath(gen.repositoryPath, gen.projectName),
	}

	for topDir, subDirs := range bp.Dirs {
		for subDir, files := range subDirs {
			dirPath := filepath.Join(projectDir, topDir, subDir)
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("failed to create dir %q: %w", dirPath, err)
			}
			for _, fileName := range files {
				filePath := filepath.Join(dirPath, fileName)
				tmplPath := "templates/" + fileName + ".tmpl"
				if err := gen.writeFile(tmplPath, filePath, data); err != nil {
					return err
				}
				if gen.verbose {
					gen.logger.Infof("created %s", filePath)
				}
			}
		}
	}

	gen.logger.Infof("next steps:")
	gen.logger.Infof("  cd %s", gen.projectName)
	gen.logger.Infof("  go mod init %s", bp.ModulePath(gen.repositoryPath, gen.projectName))
	gen.logger.Infof("  go mod tidy")
	if bp.Go != "" {
		gen.logger.Infof("  # go %s or higher required", bp.Go)
	}

	return nil
}

func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	for i, r := range name {
		if i == 0 && !unicode.IsLetter(r) {
			return fmt.Errorf("project name must start with a letter, got %q", r)
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' && r != '.' {
			return fmt.Errorf("project name contains invalid character %q", r)
		}
	}
	return nil
}

func (gen *Generator) writeFile(tmplPath, dst string, data any) error {
	if _, err := templatesFS.ReadFile(tmplPath); err == nil {
		return gen.render(tmplPath, dst, data)
	}
	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create %q: %w", dst, err)
	}
	return f.Close()
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

