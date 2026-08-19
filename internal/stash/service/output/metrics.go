package output

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	batchWriteTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "blackhole",
			Subsystem: "stash_output",
			Name:      "batch_write_total",
			Help:      "Total number of Stash output batch write executions.",
		},
		[]string{"writer", "status"},
	)

	batchWriteRecords = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "blackhole",
			Subsystem: "stash_output",
			Name:      "batch_write_records_total",
			Help:      "Total number of records in Stash output batch write executions.",
		},
		[]string{"writer", "status"},
	)

	batchWriteDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "blackhole",
			Subsystem: "stash_output",
			Name:      "batch_write_duration_seconds",
			Help:      "Duration of Stash output batch write executions.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"writer", "status"},
	)
)

func init() {
	prometheus.MustRegister(batchWriteTotal, batchWriteRecords, batchWriteDuration)
	initBatchWriteMetrics("clickhouse")
	initBatchWriteMetrics("elasticsearch")
}

func recordBatchWrite(writer, status string, records int, duration time.Duration) {
	batchWriteTotal.WithLabelValues(writer, status).Inc()
	batchWriteRecords.WithLabelValues(writer, status).Add(float64(records))
	batchWriteDuration.WithLabelValues(writer, status).Observe(duration.Seconds())
}

func initBatchWriteMetrics(writer string) {
	for _, status := range []string{"success", "failure"} {
		batchWriteTotal.WithLabelValues(writer, status)
		batchWriteRecords.WithLabelValues(writer, status)
		batchWriteDuration.WithLabelValues(writer, status)
	}
}
