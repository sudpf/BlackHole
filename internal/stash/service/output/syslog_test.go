package output

import (
	"BlackHole/internal/stash/config"
	"context"
	"testing"
)

func TestNewSyslogWriterBuildsConditionFilters(t *testing.T) {
	writer, err := NewSyslogWriter(context.Background(), &config.SyslogOutputConf{
		Conditions: [][]*config.ConditionConf{
			{
				{Key: "env", Value: "prod", Type: "match", Op: "and"},
				{Key: "service", Value: "api", Type: "contains", Op: "and"},
			},
			{
				{Key: "severity", Value: "critical", Type: "match", Op: "and"},
			},
		},
	})
	if err != nil {
		t.Fatalf("NewSyslogWriter error = %v", err)
	}

	if len(writer.Filters) != 2 {
		t.Fatalf("len(Filters) = %d, want 2", len(writer.Filters))
	}
	if !syslogFiltersMatch(writer, map[string]interface{}{"env": "prod", "service": "public-api"}) {
		t.Fatal("first condition group should match")
	}
	if !syslogFiltersMatch(writer, map[string]interface{}{"severity": "critical"}) {
		t.Fatal("second condition group should match")
	}
	if syslogFiltersMatch(writer, map[string]interface{}{"env": "dev", "service": "public-api", "severity": "warning"}) {
		t.Fatal("condition groups should not match")
	}
}

func TestNewSyslogWriterSkipsEmptyConditionGroups(t *testing.T) {
	writer, err := NewSyslogWriter(context.Background(), &config.SyslogOutputConf{
		Conditions: [][]*config.ConditionConf{
			nil,
			{nil},
			{{Key: "level", Value: "error", Type: "match", Op: "and"}},
		},
	})
	if err != nil {
		t.Fatalf("NewSyslogWriter error = %v", err)
	}

	if len(writer.Filters) != 1 {
		t.Fatalf("len(Filters) = %d, want 1", len(writer.Filters))
	}
	if !syslogFiltersMatch(writer, map[string]interface{}{"level": "error"}) {
		t.Fatal("non-empty condition group should match")
	}
}

func syslogFiltersMatch(writer *SyslogWriter, value map[string]interface{}) bool {
	for _, filter := range writer.Filters {
		if filter(value) == nil {
			return true
		}
	}
	return false
}
