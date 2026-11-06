package exporter

import (
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
)

// gatedCollector only exports metrics while its device is reachable, so an
// offline device stops reporting instead of repeating its last values. That
// lets missing-data alerts fire.
type gatedCollector struct {
	prometheus.Collector
	online *atomic.Bool
}

func (g gatedCollector) Collect(ch chan<- prometheus.Metric) {
	if g.online.Load() {
		g.Collector.Collect(ch)
	}
}
