```go
package main

import (
 "github.com/prometheus/client_golang/prometheus"
 "log"
 "net/http"

 "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
 gauge := prometheus.NewGauge(prometheus.GaugeOpts{
  Namespace: "gNamespace", Subsystem: "gSubsystem", Name: "gName", Help: "gHelp",
  ConstLabels: prometheus.Labels{"gKey1": "gValue1", "gKey2": "gValue2"}})

 counter := prometheus.NewCounter(prometheus.CounterOpts{
  Namespace: "cNamespace", Subsystem: "cSubsystem", Name: "cName", Help: "cHelp",
  ConstLabels: prometheus.Labels{"cKey1": "cValue1", "cKey2": "cValue2"}})

 summary := prometheus.NewSummary(prometheus.SummaryOpts{
  Namespace: "sNamespace", Subsystem: "sSubsystem", Name: "sName", Help: "sHelp",
  ConstLabels: prometheus.Labels{"sKey1": "sValue1", "sKey2": "sValue2"}})

 histogram := prometheus.NewHistogram(prometheus.HistogramOpts{
  Namespace: "hNamespace", Subsystem: "hSubsystem", Name: "hName", Help: "hHelp",
  ConstLabels: prometheus.Labels{"hKey1": "hValue1", "hKey2": "hValue2"}, Buckets: []float64{}})

 registry := prometheus.NewRegistry()
 registry.MustRegister(gauge, counter, summary, histogram)

 http.Handle("/update", myHandler{gauge, counter, summary, histogram})

 http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
 log.Fatal(http.ListenAndServe(":8080", nil))
}

type myHandler struct {
 gauge     prometheus.Gauge
 counter   prometheus.Counter
 summary   prometheus.Summary
 histogram prometheus.Histogram
}

func (handler myHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
 handler.gauge.Set(123)
 handler.counter.Inc()
 handler.summary.Observe(10)
 handler.histogram.Observe(30)
}
```
