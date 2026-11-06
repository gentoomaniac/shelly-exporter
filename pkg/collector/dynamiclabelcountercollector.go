package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type DynamicLabelCounterCollectorOpts struct {
	Namespace     string
	Name          string
	Help          string
	DynamicLabels []string
	ConstLabels   prometheus.Labels
}

func NewDynamicLabelCounterCollector(opts DynamicLabelCounterCollectorOpts, valueFunc func() float64, labelValueFunc func() []string) *DynamicLabelCounterCollector {
	desc := prometheus.NewDesc(
		fmt.Sprintf("%s_%s", opts.Namespace, opts.Name),
		opts.Help,
		opts.DynamicLabels,
		opts.ConstLabels,
	)

	c := &DynamicLabelCounterCollector{
		description:       desc,
		dynamicLabels:     opts.DynamicLabels,
		dynamicLabelsFunc: labelValueFunc,
		valueFunc:         valueFunc,
	}

	return c
}

type DynamicLabelCounterCollector struct {
	description       *prometheus.Desc
	dynamicLabels     []string
	dynamicLabelsFunc func() []string

	valueFunc func() float64
}

func (c *DynamicLabelCounterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.description
}

func (c *DynamicLabelCounterCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.description,
		prometheus.CounterValue,
		c.valueFunc(),
		c.dynamicLabelsFunc()...,
	)
}
