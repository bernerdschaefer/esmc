package esmc

import (
	"github.com/prometheus/client_golang"
	"github.com/prometheus/client_golang/metrics"
	"time"
)

var (
	requestCount     = metrics.NewCounter()
	requestDuration  = metrics.NewCounter()
	requestDurations = metrics.NewDefaultHistogram()

	reportedRequestDuration  = metrics.NewCounter()
	reportedRequestDurations = metrics.NewDefaultHistogram()
)

func updateRequestMetrics(labels map[string]string, duration time.Duration) {
	requestCount.Increment(labels)
	requestDuration.IncrementBy(labels, float64(duration))
	requestDurations.Add(labels, float64(duration))
}

func updateReportedMetrics(labels map[string]string, took int) {
	duration := float64(time.Duration(took) * time.Millisecond)

	reportedRequestDuration.IncrementBy(labels, duration)
	reportedRequestDurations.Add(labels, duration)
}

func init() {
	registry.Register(
		"elasticsearch_client_requests",
		"A counter of the total number of requests to an ES cluster",
		registry.NilLabels,
		requestCount,
	)
	registry.Register(
		"elasticsearch_client_request_total_duration_nanoseconds",
		"The total amount of time spent executing requests (nanoseconds)",
		registry.NilLabels,
		requestDuration,
	)
	registry.Register(
		"elasticsearch_client_request_durations_nanoseconds",
		"The amounts of time spent executing requests (nanoseconds)",
		registry.NilLabels,
		requestDurations,
	)
	registry.Register(
		"elasticsearch_client_reported_request_total_duration_nanoseconds",
		"The total amount of time spent executing requests as reported by elasticsearch (nanoseconds)",
		registry.NilLabels,
		reportedRequestDuration,
	)
	registry.Register(
		"elasticsearch_client_reported_request_durations_nanoseconds",
		"The amounts of time spent executing requests as reported by elasticsearch (nanoseconds)",
		registry.NilLabels,
		reportedRequestDurations,
	)
}
