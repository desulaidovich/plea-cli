// Package blueprint reads and parses the blueprint.yaml file that defines
// the directory and file structure for a generated project.
package blueprint

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const RootFile = "blueprint.yaml"

type Blueprint struct {
	Go     string                         `yaml:"go"`
	Module string                         `yaml:"module"`
	Dirs   map[string]map[string][]string `yaml:"dirs"`
}

// ModulePath returns the module path with {{repositoryPath}} and {{serviceName}} replaced.
func (bp *Blueprint) ModulePath(repositoryPath, serviceName string) string {
	r := strings.NewReplacer(
		"{{repositoryPath}}", repositoryPath,
		"{{serviceName}}", serviceName,
	)
	return r.Replace(bp.Module)
}

func Read() (*Blueprint, error) {
	raw, err := os.ReadFile("./" + RootFile)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", RootFile, err)
	}

	var bp Blueprint
	if err = yaml.Unmarshal(raw, &bp); err != nil {
		return nil, fmt.Errorf("parse %s: %w", RootFile, err)
	}

	return &bp, nil
}
