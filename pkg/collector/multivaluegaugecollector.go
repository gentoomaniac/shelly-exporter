package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

// A simple struct to hold the "row" of data for one sensor
type MetricPoint struct {
	Value       float64
	LabelValues []string
}

type MultiValueGaugeCollector struct {
	description *prometheus.Desc
	pointsFunc  func() []MetricPoint
}

type MultiValueGaugeCollectorOpts struct {
	Namespace     string
	Name          string
	Help          string
	DynamicLabels []string
	ConstLabels   prometheus.Labels
}

func NewMultiValueGaugeCollector(opts MultiValueGaugeCollectorOpts, pointsFunc func() []MetricPoint) *MultiValueGaugeCollector {
	desc := prometheus.NewDesc(
		prometheus.BuildFQName(opts.Namespace, "", opts.Name),
		opts.Help,
		opts.DynamicLabels,
		opts.ConstLabels,
	)
	return &MultiValueGaugeCollector{
		description: desc,
		pointsFunc:  pointsFunc,
	}
}

func (c *MultiValueGaugeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.description
}

func (c *MultiValueGaugeCollector) Collect(ch chan<- prometheus.Metric) {
	points := c.pointsFunc()
	for _, p := range points {
		ch <- prometheus.MustNewConstMetric(
			c.description,
			prometheus.GaugeValue,
			p.Value,
			p.LabelValues...,
		)
	}
}
