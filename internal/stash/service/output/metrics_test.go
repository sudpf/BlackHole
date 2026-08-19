package output

import (
	"testing"
	"time"

	dto "github.com/prometheus/client_model/go"
)

func TestRecordBatchWrite(t *testing.T) {
	beforeTotal := counterValue(t, batchWriteTotal.WithLabelValues("test-writer", "success"))
	beforeRecords := counterValue(t, batchWriteRecords.WithLabelValues("test-writer", "success"))

	recordBatchWrite("test-writer", "success", 3, 10*time.Millisecond)

	if got := counterValue(t, batchWriteTotal.WithLabelValues("test-writer", "success")); got != beforeTotal+1 {
		t.Fatalf("batch write total = %v, want %v", got, beforeTotal+1)
	}
	if got := counterValue(t, batchWriteRecords.WithLabelValues("test-writer", "success")); got != beforeRecords+3 {
		t.Fatalf("batch write records = %v, want %v", got, beforeRecords+3)
	}
}

func counterValue(t *testing.T, metric interface{ Write(*dto.Metric) error }) float64 {
	t.Helper()

	var m dto.Metric
	if err := metric.Write(&m); err != nil {
		t.Fatalf("write metric: %v", err)
	}
	if m.Counter == nil || m.Counter.Value == nil {
		return 0
	}
	return *m.Counter.Value
}
