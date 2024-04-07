package reporter

// WavefrontMetricsReporter report go-metrics to wavefront
type Reporter interface {
	// Starts reporting metrics to Wavefront at a given interval.
	Start()

	// Stops reporting metrics and closes the reporter.
	Close()

	// Reports the metrics to Wavefront just once. Can be used to manually report metrics to Wavefront outside of Start.
	Report()

	// Gets the count of errors in reporting metrics to Wavefront.
	//ErrorsCount() int64

	// RegisterMetric register the given metric under the given name and tags
	// return RegistryError if the metric is not registered
	//RegisterMetric(name string, metric interface{}, tags map[string]string) error

	// GetMetric get the metric by the given name and tags or nil if none is registered.
	//GetMetric(name string, tags map[string]string) interface{}

	// GetOrRegisterMetric gets an existing metric or registers the given one.
	// The interface can be the metric to register if not found in registry,
	// or a function returning the metric for lazy instantiation.
	//GetOrRegisterMetric(name string, i interface{}, tags map[string]string) interface{}

	// UnregisterMetric Unregister the metric with the given name.
	//UnregisterMetric(name string, tags map[string]string)
}
