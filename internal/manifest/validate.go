package manifest

import (
	"fmt"
	"path/filepath"
)

var validNames = map[string]struct{}{
	"plea.yaml": {},
	"plea.yml":  {},
	"plea.json": {},
}

func validateFileName(path string) (string, error) {
	name := filepath.Base(path)
	if _, ok := validNames[name]; !ok {
		return "", fmt.Errorf("invalid file name %q: expected plea.yaml, plea.yml or plea.json", name)
	}
	return path, nil
}
