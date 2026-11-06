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
	// ReadyFunc is optional. When it returns false, Collect emits nothing,
	// e.g. while the device data has not been read yet.
	ReadyFunc func() bool
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
		readyFunc:         opts.ReadyFunc,
	}

	return c
}

type DynamicLabelCounterCollector struct {
	description       *prometheus.Desc
	dynamicLabels     []string
	dynamicLabelsFunc func() []string

	valueFunc func() float64
	readyFunc func() bool
}

func (c *DynamicLabelCounterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.description
}

func (c *DynamicLabelCounterCollector) Collect(ch chan<- prometheus.Metric) {
	if c.readyFunc != nil && !c.readyFunc() {
		return
	}
	ch <- prometheus.MustNewConstMetric(
		c.description,
		prometheus.CounterValue,
		c.valueFunc(),
		c.dynamicLabelsFunc()...,
	)
}
