package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveConfigFileUsesWorkingDirectory(t *testing.T) {
	workingDir := t.TempDir()
	oldWorkingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	t.Cleanup(func() {
		if chdirErr := os.Chdir(oldWorkingDir); chdirErr != nil {
			t.Fatalf("restore working directory: %v", chdirErr)
		}
	})

	if chdirErr := os.Chdir(workingDir); chdirErr != nil {
		t.Fatalf("change working directory: %v", chdirErr)
	}

	got, err := resolveConfigFile(filepath.Join("conf", "stash.yaml"))
	if err != nil {
		t.Fatalf("resolve config file: %v", err)
	}

	want := filepath.Join(workingDir, "conf", "stash.yaml")
	if got != want {
		t.Fatalf("config file = %q, want %q", got, want)
	}
}

func TestResolveConfigFileRejectsEmptyPath(t *testing.T) {
	if _, err := resolveConfigFile(" "); err == nil {
		t.Fatal("expected empty config file to be rejected")
	}
}
