package esmc

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var (
	requestCount     = prometheus.NewCounter()
	requestDuration  = prometheus.NewCounter()
	requestDurations = prometheus.NewDefaultHistogram()

	reportedRequestDuration  = prometheus.NewCounter()
	reportedRequestDurations = prometheus.NewDefaultHistogram()
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
	prometheus.Register(
		"esmc_requests",
		"A counter of the total number of requests to an ES cluster",
		prometheus.NilLabels,
		requestCount,
	)
	prometheus.Register(
		"esmc_request_total_duration_nanoseconds",
		"The total amount of time spent executing requests (nanoseconds)",
		prometheus.NilLabels,
		requestDuration,
	)
	prometheus.Register(
		"esmc_request_durations_nanoseconds",
		"The amounts of time spent executing requests (nanoseconds)",
		prometheus.NilLabels,
		requestDurations,
	)
	prometheus.Register(
		"esmc_reported_request_total_duration_nanoseconds",
		"The total amount of time spent executing requests as reported by elasticsearch (nanoseconds)",
		prometheus.NilLabels,
		reportedRequestDuration,
	)
	prometheus.Register(
		"esmc_reported_request_durations_nanoseconds",
		"The amounts of time spent executing requests as reported by elasticsearch (nanoseconds)",
		prometheus.NilLabels,
		reportedRequestDurations,
	)
}
