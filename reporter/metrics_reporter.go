package reporter

import (
	"github.com/Edward-Alphonse/saywo_metrics/sender"
	"github.com/Edward-Alphonse/saywo_metrics/utils"
	"github.com/rcrowley/go-metrics"
	"strconv"
	"strings"
	"sync"
	"time"
)

var once sync.Once

type MetricsReporter struct {
	sender sender.MetricSender
	//application   application.Tags
	source       string
	prefix       string
	addSuffix    bool
	interval     time.Duration
	ticker       *time.Ticker
	percentiles  []float64     // Percentiles to export from timers and histograms
	durationUnit time.Duration // Time conversion unit for durations
	//metrics      map[string]interface{} // for Wavefron specific metrics tyoes, like Histograms
	errors chan error
	start  chan bool
	close  chan bool
	//errorsCount  int64
	//errorDebug   bool
	//autoStart     bool
	mux           sync.Mutex
	registry      metrics.Registry
	runtimeMetric bool // for getting the go runtime metrics
}

// Option allows WavefrontReporter customization
type Option func(*MetricsReporter)

// LogErrors tag for metrics
//func LogErrors(debug bool) Option {
//	return func(args *reporter) {
//		args.errorDebug = debug
//	}
//}
//
//// Source tag for metrics
//func Source(source string) Option {
//	return func(args *reporter) {
//		args.source = source
//	}
//}

// Interval change the default 1 minute reporting interval.
func Interval(interval time.Duration) Option {
	return func(args *MetricsReporter) {
		args.interval = interval
	}
}

//// DisableAutoStart prevents the Reporter from automatically reporting when created.
//func DisableAutoStart() Option {
//	return func(args *MetricsReporter) {
//		args.autoStart = false
//	}
//}

// Prefix for the metrics name
func Prefix(prefix string) Option {
	return func(args *MetricsReporter) {
		args.prefix = strings.TrimSuffix(prefix, ".")
	}
}

// AddSuffix adds a metric suffix based on the metric type ('.count', '.value')
func DisableAddSuffix(addSuffix bool) Option {
	return func(args *MetricsReporter) {
		args.addSuffix = addSuffix
	}
}

// CustomRegistry allows overriding the registry used by the reporter.
//func CustomRegistry(registry metrics.Registry) Option {
//	return func(args *reporter) {
//		args.registry = registry
//	}
//}

// Enable Runtime Metric collection
func EnableRuntimeMetric(enable bool) Option {
	return func(args *MetricsReporter) {
		args.runtimeMetric = enable
	}
}

// Application tag for the metrics
//func ApplicationTag(application application.Tags) Option {
//	return func(args *reporter) {
//		args.application = application
//	}
//}

// NewMetricsReporter create a WavefrontMetricsReporter
func NewMetricsReporter(sender sender.MetricSender, setters ...Option) *MetricsReporter {
	r := &MetricsReporter{
		sender: sender,
		//source:        hostname(),
		interval:     time.Second * 5,
		percentiles:  []float64{0.5, 0.75, 0.95, 0.99, 0.999},
		durationUnit: time.Nanosecond,
		//metrics:      make(map[string]interface{}),
		addSuffix: true,
		//errorsCount:   0,
		//autoStart:     true,
		runtimeMetric: false,
	}

	for _, setter := range setters {
		setter(r)
	}

	//if r.registry == nil {
	r.registry = metrics.DefaultRegistry
	//}

	if r.runtimeMetric == true {
		metrics.RegisterRuntimeMemStats(r.registry)
	}

	r.ticker = time.NewTicker(r.interval)
	r.close = make(chan bool, 1)
	r.start = make(chan bool, 1)
	r.errors = make(chan error)
	//go r.process()
	//if r.autoStart {
	//	r.Start()
	//}
	return r
}

//func ()

// Deprecated: use NewMetricsReporter
//func NewReporter(sender wf.Sender, application application.Tags, setters ...Option) Reporter {
//	return NewMetricsReporter(sender, append(setters, ApplicationTag(application))...)
//}

//func (r *reporter) ErrorsCount() int64 {
//	return atomic.LoadInt64(&r.errorsCount)
//}

//func aaa() {
//	once := sync.Once{}
//	once.Do(func() {
//		go
//	})
//}

func (r *MetricsReporter) runOnce() {
	once.Do(func() {
		go r.run()
	})
}

func (r *MetricsReporter) run() {
	running := false
	for {
		select {
		case <-r.ticker.C:
			if running {
				go r.Report()
			}
		case error := <-r.errors:
			if error != nil {
				//atomic.AddInt64(&r.errorsCount, 1)
				//if r.errorDebug {
				//	log.Printf("reporter error: %v\n", error)
				//}
			}
		case <-r.start:
			r.sender.Start()
			running = true
		case <-r.close:
			r.Report()
			//r.sender.Close()
			return
		}
	}
}

func (r *MetricsReporter) Report() {
	r.mux.Lock()
	defer r.mux.Unlock()
	registry := r.registry
	//lastErrorsCount := r.ErrorsCount()

	if r.runtimeMetric == true {
		metrics.CaptureRuntimeMemStatsOnce(registry)
	}

	registry.Each(func(key string, metric interface{}) {
		name, tags := utils.DecodeKey(key)

		//for t, v := range r.application.Map() {
		//	if _, ok := tags[t]; !ok && len(v) > 0 {
		//		tags[t] = v
		//	}
		//}
		timestamp := time.Now().Unix()
		switch metric.(type) {
		case metrics.Counter:
			//if hasDeltaPrefix(name) {
			//	r.reportDelta(name, metric.(metrics.Counter), tags)
			//} else {
			r.errors <- r.sender.SendMetric(r.prepareName(name, "count"), tags, float64(metric.(metrics.Counter).Count()), timestamp, r.source)
			//}
		case metrics.Gauge:
			r.errors <- r.sender.SendMetric(r.prepareName(name, "value"), tags, float64(metric.(metrics.Gauge).Value()), timestamp, r.source)
		case metrics.GaugeFloat64:
			r.errors <- r.sender.SendMetric(r.prepareName(name, "value"), tags, float64(metric.(metrics.GaugeFloat64).Value()), timestamp, r.source)
		//case Histogram:
		//	r.reportWFHistogram(name, metric.(Histogram), tags)
		case metrics.Histogram:
			r.reportHistogram(name, metric.(metrics.Histogram), tags, timestamp)
		case metrics.Meter:
			r.reportMeter(name, metric.(metrics.Meter), tags, timestamp)
		case metrics.Timer:
			r.reportTimer(name, metric.(metrics.Timer), tags, timestamp)
		}
	})
	//actualErrorsCount := r.ErrorsCount()
	//if actualErrorsCount != lastErrorsCount {
	//	log.Printf("!!! There was %d errors on the last reporting cycle !!!", (actualErrorsCount - lastErrorsCount))
	//}
}

//func (r *MetricsReporter) reportDelta(name string, metric metrics.Counter, tags map[string]string) {
//	var prunedName string
//	if strings.HasPrefix(name, deltaPrefix) {
//		prunedName = name[deltaPrefixSize:]
//	} else if strings.HasPrefix(name, altDeltaPrefix) {
//		prunedName = name[altDeltaPrefixSize:]
//	}
//	value := metric.Count()
//	metric.Dec(value)
//
//	r.errors <- r.sender.SendDeltaCounter(deltaPrefix+r.prepareName(prunedName, "count"), tags, float64(value), r.source)
//}
//
//func (r *MetricsReporter) reportWFHistogram(metricName string, h Histogram, tags map[string]string) {
//	distributions := h.Distributions()
//	hgs := map[histogram.Granularity]bool{h.Granularity(): true}
//	for _, distribution := range distributions {
//		if len(distribution.Centroids) > 0 {
//			r.errors <- r.sender.SendDistribution(r.prepareName(metricName), tags, distribution.Centroids, hgs, distribution.Timestamp.Unix(), r.source)
//		}
//	}
//}

func (r *MetricsReporter) reportHistogram(name string, metric metrics.Histogram, tags map[string]string, timestamp int64) {
	h := metric.Snapshot()
	ps := h.Percentiles(r.percentiles)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".count"), tags, float64(h.Count()), timestamp, r.source)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".min"), tags, float64(h.Min()), timestamp, r.source)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".max"), tags, float64(h.Max()), timestamp, r.source)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".mean"), tags, h.Mean(), timestamp, r.source)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".std-dev"), tags, h.StdDev(), timestamp, r.source)
	for psIdx, psKey := range r.percentiles {
		key := strings.Replace(strconv.FormatFloat(psKey*100.0, 'f', -1, 64), ".", "", 1)
		r.errors <- r.sender.SendMetric(r.prepareName(name+"."+key+"-percentile"), tags, ps[psIdx], timestamp, r.source)
	}
}

func (r *MetricsReporter) reportMeter(name string, metric metrics.Meter, tags map[string]string, timestamp int64) {
	m := metric.Snapshot()
	r.errors <- r.sender.SendMetric(r.prepareName(name+".count"), tags, float64(m.Count()), timestamp, r.source)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".one-minute"), tags, m.Rate1(), timestamp, r.source)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".five-minute"), tags, m.Rate5(), timestamp, r.source)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".fifteen-minute"), tags, m.Rate15(), timestamp, r.source)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".mean"), tags, m.RateMean(), timestamp, r.source)
}

func (r *MetricsReporter) reportTimer(name string, metric metrics.Timer, tags map[string]string, timestamp int64) {
	t := metric.Snapshot()
	du := float64(r.durationUnit)
	ps := t.Percentiles(r.percentiles)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".count"), tags, float64(t.Count()), timestamp, r.source)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".min"), tags, float64(t.Min()/int64(du)), timestamp, r.source)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".max"), tags, float64(t.Max()/int64(du)), timestamp, r.source)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".mean"), tags, t.Mean()/du, timestamp, r.source)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".std-dev"), tags, t.StdDev()/du, timestamp, r.source)
	for psIdx, psKey := range r.percentiles {
		key := strings.Replace(strconv.FormatFloat(psKey*100.0, 'f', -1, 64), ".", "", 1)
		r.errors <- r.sender.SendMetric(r.prepareName(name+"."+key+"-percentile"), tags, ps[psIdx]/du, timestamp, r.source)
	}
	r.errors <- r.sender.SendMetric(r.prepareName(name+".one-minute"), tags, t.Rate1(), timestamp, r.source)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".five-minute"), tags, t.Rate5(), timestamp, r.source)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".fifteen-minute"), tags, t.Rate15(), timestamp, r.source)
	r.errors <- r.sender.SendMetric(r.prepareName(name+".mean-rate"), tags, t.RateMean(), timestamp, r.source)
}

func (r *MetricsReporter) prepareName(name string, suffix ...string) string {
	if len(r.prefix) > 0 {
		name = r.prefix + "." + name
	}

	if r.addSuffix {
		for _, s := range suffix {
			name += "." + s
		}
	}

	return name
}

func (r *MetricsReporter) Start() {
	r.runOnce()
	r.start <- true
}

func (r *MetricsReporter) Close() {
	r.close <- true
}

// RegistryError returned if there is any error on RegisterMetric
//type RegistryError string
//
//func (err RegistryError) Error() string {
//	return fmt.Sprintf("Registry Error: %s", string(err))
//}

// RegisterMetric register the given metric under the given name and tags
// return RegistryError if the metric is not registered
//func (r *MetricsReporter) RegisterMetric(name string, metric interface{}, tags map[string]string) error {
//	key := EncodeKey(name, tags)
//	err := r.registry.Register(key, metric)
//	if err != nil {
//		return err
//	}
//	m := r.GetMetric(name, tags)
//	if m == nil {
//		return RegistryError(fmt.Sprintf("Metric '%s'(%s) not registered.", name, reflect.TypeOf(metric).String()))
//	}
//	return nil
//}

// GetMetric get the metric by the given name and tags or nil if none is registered.
//func (r *MetricsReporter) GetMetric(name string, tags map[string]string) interface{} {
//	key := EncodeKey(name, tags)
//	return r.registry.Get(key)
//}

// GetOrRegisterMetric gets an existing metric or registers the given one.
// The interface can be the metric to register if not found in registry,
// or a function returning the metric for lazy instantiation.
//func (r *MetricsReporter) GetOrRegisterMetric(name string, i interface{}, tags map[string]string) interface{} {
//	key := EncodeKey(name, tags)
//	return r.registry.GetOrRegister(key, i)
//}

// UnregisterMetric Unregister the metric with the given name.
//func (r *MetricsReporter) UnregisterMetric(name string, tags map[string]string) {
//	key := utils.EncodeKey(name, tags)
//	r.registry.Unregister(key)
//}

//func hostname() string {
//	name, err := os.Hostname()
//	if err != nil {
//		name = "go-metrics-wavefront"
//	}
//	return name
//}
