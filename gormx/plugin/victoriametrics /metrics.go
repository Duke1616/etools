package victoriametrics

//// Register various metrics.
//// Metric name may contain labels in Prometheus format - see below.
//var (
//	// Register counter without labels.
//	requestsTotal = metrics.NewCounter("requests_total")
//
//	// Register summary with a single label.
//	requestDuration = metrics.NewSummary(`requests_duration_seconds{path="/foobar/baz"}`)
//
//	// Register gauge with two labels.
//	queueSize = metrics.NewGauge(`queue_size{queue="foobar",topic="baz"}`, func() float64 {
//		return float64(foobarQueue.Len())
//	})
//
//	// Register histogram with a single label.
//	responseSize = metrics.NewHistogram(`response_size{path="/foo/bar"}`)
//)
//
//func NewSummaryVec() {
//	metrics.NewSummaryExt()
//	metrics.PushMetrics()
//	prometheus.NewSummaryVec()
//}
