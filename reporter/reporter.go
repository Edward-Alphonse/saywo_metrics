package reporter

// WavefrontMetricsReporter report go-metrics to wavefront
type Reporter interface {
	// Starts reporting metrics to Wavefront at a given interval.
	Start()

	// Stops reporting metrics and closes the reporter.
	Close()

	// Reports the metrics to Wavefront just once. Can be used to manually report metrics to Wavefront outside of Start.
	Report()
}
