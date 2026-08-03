package configpath

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveUsesWorkingDirectory(t *testing.T) {
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

	got, err := Resolve(filepath.Join("conf", "app.yaml"))
	if err != nil {
		t.Fatalf("resolve config file: %v", err)
	}

	want := filepath.Join(workingDir, "conf", "app.yaml")
	if got != want {
		t.Fatalf("config file = %q, want %q", got, want)
	}
}

func TestResolveCleansAbsolutePath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "conf", "..", "app.yaml")
	got, err := Resolve(path)
	if err != nil {
		t.Fatalf("resolve config file: %v", err)
	}

	want := filepath.Clean(path)
	if got != want {
		t.Fatalf("config file = %q, want %q", got, want)
	}
}

func TestResolveRejectsEmptyPath(t *testing.T) {
	if _, err := Resolve(" "); err == nil {
		t.Fatal("expected empty config file to be rejected")
	}
}
