package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type DynamicLabelGaugeCollectorOpts struct {
	Namespace     string
	Name          string
	Help          string
	DynamicLabels []string
	ConstLabels   prometheus.Labels
}

func NewDynamicLabelGaugeCollector(opts DynamicLabelGaugeCollectorOpts, valueFunc func() float64, labelValueFunc func() []string) *DynamicLabelGaugeCollector {
	desc := prometheus.NewDesc(
		fmt.Sprintf("%s:%s:gauge", opts.Namespace, opts.Name),
		opts.Help,
		opts.DynamicLabels,
		opts.ConstLabels,
	)

	c := &DynamicLabelGaugeCollector{
		description:       desc,
		dynamicLabels:     opts.DynamicLabels,
		dynamicLabelsFunc: labelValueFunc,
		valueFunc:         valueFunc,
	}

	return c
}

type DynamicLabelGaugeCollector struct {
	description       *prometheus.Desc
	dynamicLabels     []string
	dynamicLabelsFunc func() []string

	valueFunc func() float64
}

func (c *DynamicLabelGaugeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.description
}

func (c *DynamicLabelGaugeCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.description,
		prometheus.GaugeValue,
		c.valueFunc(), // value
		c.dynamicLabelsFunc()...,
	)
}
