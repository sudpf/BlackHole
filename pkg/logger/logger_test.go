package logger

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRotatingWriterRejectsLessThanOneMB(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "test.log")
	if _, err := RotatingWriter(filename, "10240"); err == nil {
		t.Fatal(`RotatingWriter("test.log", "10240") expected error`)
	}
}

func TestRotatingWriterReturnsPrepareLogFileError(t *testing.T) {
	dir := t.TempDir()
	notDir := filepath.Join(dir, "not-dir")
	if err := os.WriteFile(notDir, []byte("file"), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	if _, err := RotatingWriter(filepath.Join(notDir, "test.log"), "1m"); err == nil {
		t.Fatal("RotatingWriter expected prepare log file error")
	}
}
