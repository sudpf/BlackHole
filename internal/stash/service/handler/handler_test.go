package handler

import (
	"context"
	"errors"
	"testing"
)

type testWriter struct {
	writeErr error
	closed   bool
}

func (w *testWriter) Write(map[string]interface{}) error {
	return w.writeErr
}

func (w *testWriter) Close() error {
	w.closed = true
	return nil
}

func TestConsumeReturnsWriterErrors(t *testing.T) {
	failingWriter := &testWriter{writeErr: errors.New("write failed")}
	handler := NewHandler()
	handler.AddWriters(&testWriter{}, failingWriter)

	err := handler.Consume(context.Background(), "", `{"message":"test"}`)
	if !errors.Is(err, failingWriter.writeErr) {
		t.Fatalf("Consume() error = %v, want wrapped %v", err, failingWriter.writeErr)
	}
}

func TestCloseClosesAllWriters(t *testing.T) {
	first := &testWriter{}
	second := &testWriter{}
	handler := NewHandler()
	handler.AddWriters(first, second)

	if err := handler.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if !first.closed || !second.closed {
		t.Fatal("Close() did not close every writer")
	}
}
