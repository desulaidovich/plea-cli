// Package manifest reads and parses plea.yaml/plea.json configuration files.
package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ReadFile(path string) (*Config, error) {
	fileName, err := validateFileName(path)
	if err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", fileName, err)
	}

	var cfg Config
	switch filepath.Ext(fileName) {
	case ".json":
		err = json.Unmarshal(raw, &cfg)
	default:
		err = yaml.Unmarshal(raw, &cfg)
	}
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", fileName, err)
	}

	return &cfg, nil
}
