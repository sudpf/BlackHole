package configpath

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Resolve(file string) (string, error) {
	if strings.TrimSpace(file) == "" {
		return "", fmt.Errorf("config file is required")
	}
	if filepath.IsAbs(file) {
		return filepath.Clean(file), nil
	}

	absPath, err := filepath.Abs(file)
	if err != nil {
		return "", fmt.Errorf("resolve config file %q: %w", file, err)
	}
	return absPath, nil
}

func ExecutableName() (string, error) {
	appPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("get executable path: %w", err)
	}
	return filepath.Base(appPath), nil
}
