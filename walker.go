package generator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dohernandez/errors"
)

// FindModelSourcePath finds the source file of the model.
func FindModelSourcePath(model, dir string) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") || strings.HasPrefix(file.Name(), "_") {
			continue
		}

		if file.IsDir() {
			continue
		}

		path := filepath.Join(dir, file.Name())

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			continue
		}

		content, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return "", err
		}

		if strings.Contains(string(content), "type "+model+" struct") {
			return path, nil
		}
	}

	return "", errors.New("model source file not found")
}
