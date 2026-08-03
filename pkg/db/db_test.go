package db

import (
	"BlackHole/pkg/requestctx"
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type queryTestRecord struct {
	ID   int
	Name string
}

func TestLogrusAdapterIncludesTraceID(t *testing.T) {
	traceID := "4bf92f3577b34da6a3ce929d0e0e4736"
	var output bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&output)
	logger.SetFormatter(&CustomFormatter{})
	adapter := NewLogrusAdapter(logger)

	ctx := requestctx.WithScope(context.Background(), requestctx.Scope{TraceID: traceID})
	adapter.Trace(ctx, time.Now(), func() (string, int64) {
		return "SELECT 1", 1
	}, nil)

	var fields map[string]interface{}
	if err := json.Unmarshal(output.Bytes(), &fields); err != nil {
		t.Fatalf("unmarshal log: %v: %s", err, output.String())
	}
	if fields["trace_id"] != traceID {
		t.Fatalf("trace_id = %v, want %s", fields["trace_id"], traceID)
	}
	if _, ok := fields["db"]; ok {
		t.Fatalf("db should not be present: %v", fields)
	}
	if _, ok := fields["elapsed_ms"]; !ok {
		t.Fatalf("elapsed_ms missing: %v", fields)
	}
}

func TestLogrusAdapterMarksStartupWhenTraceIDMissing(t *testing.T) {
	var output bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&output)
	logger.SetFormatter(&CustomFormatter{})
	adapter := NewLogrusAdapter(logger)

	adapter.Trace(context.Background(), time.Now(), func() (string, int64) {
		return "SELECT 1", 1
	}, nil)

	var fields map[string]interface{}
	if err := json.Unmarshal(output.Bytes(), &fields); err != nil {
		t.Fatalf("unmarshal log: %v: %s", err, output.String())
	}
	if fields["trace_id"] != "system" {
		t.Fatalf("trace_id = %v, want system", fields["trace_id"])
	}
}

func TestQueryUsesOptionsWithoutMutatingConditions(t *testing.T) {
	gormDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := gormDB.AutoMigrate(&queryTestRecord{}); err != nil {
		t.Fatalf("migrate records: %v", err)
	}
	if err := gormDB.Create([]queryTestRecord{
		{ID: 1, Name: "alice"},
		{ID: 2, Name: "alice"},
		{ID: 3, Name: "bob"},
	}).Error; err != nil {
		t.Fatalf("seed records: %v", err)
	}

	database := &SQLiteDatabase{DB: gormDB}
	conditions := map[string]interface{}{
		"name": "alice",
	}

	var records []queryTestRecord
	if _, err := database.Query(context.Background(), &records, conditions, &QueryOptions{
		PageNo:      1,
		PageSize:    1,
		OrderColumn: "id",
		Order:       OrderDesc,
	}); err != nil {
		t.Fatalf("query records: %v", err)
	}

	if len(records) != 1 || records[0].ID != 2 {
		t.Fatalf("records = %+v, want only id 2", records)
	}
	if len(conditions) != 1 || conditions["name"] != "alice" {
		t.Fatalf("conditions mutated: %+v", conditions)
	}
}

func TestQueryExUsesOptions(t *testing.T) {
	gormDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := gormDB.AutoMigrate(&queryTestRecord{}); err != nil {
		t.Fatalf("migrate records: %v", err)
	}
	if err := gormDB.Create([]queryTestRecord{
		{ID: 1, Name: "alice"},
		{ID: 2, Name: "bob"},
	}).Error; err != nil {
		t.Fatalf("seed records: %v", err)
	}

	database := &SQLiteDatabase{DB: gormDB}

	var records []queryTestRecord
	if _, err := database.QueryEx(context.Background(), &records, queryTestRecord{}, &QueryOptions{
		PageNo:      1,
		PageSize:    1,
		OrderColumn: "id",
		Order:       OrderDesc,
	}); err != nil {
		t.Fatalf("query records: %v", err)
	}

	if len(records) != 1 || records[0].ID != 2 {
		t.Fatalf("records = %+v, want only id 2", records)
	}
}

func TestQueryAllowsNilOptions(t *testing.T) {
	gormDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := gormDB.AutoMigrate(&queryTestRecord{}); err != nil {
		t.Fatalf("migrate records: %v", err)
	}
	if err := gormDB.Create([]queryTestRecord{
		{ID: 1, Name: "alice"},
		{ID: 2, Name: "bob"},
	}).Error; err != nil {
		t.Fatalf("seed records: %v", err)
	}

	database := &SQLiteDatabase{DB: gormDB}

	var records []queryTestRecord
	if _, err := database.Query(context.Background(), &records, map[string]interface{}{}, nil); err != nil {
		t.Fatalf("query records: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("records length = %d, want 2", len(records))
	}
}
