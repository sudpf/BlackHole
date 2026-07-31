package logger

import (
	"path/filepath"
	"testing"
)

func TestRotatingWriterRejectsLessThanOneMB(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "test.log")
	if _, err := RotatingWriter(filename, "10240"); err == nil {
		t.Fatal(`RotatingWriter("test.log", "10240") expected error`)
	}
}
