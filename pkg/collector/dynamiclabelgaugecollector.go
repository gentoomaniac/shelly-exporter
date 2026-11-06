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
	// ReadyFunc is optional. When it returns false, Collect emits nothing,
	// e.g. while the device data has not been read yet.
	ReadyFunc func() bool
}

func NewDynamicLabelGaugeCollector(opts DynamicLabelGaugeCollectorOpts, valueFunc func() float64, labelValueFunc func() []string) *DynamicLabelGaugeCollector {
	desc := prometheus.NewDesc(
		fmt.Sprintf("%s_%s", opts.Namespace, opts.Name),
		opts.Help,
		opts.DynamicLabels,
		opts.ConstLabels,
	)

	c := &DynamicLabelGaugeCollector{
		description:       desc,
		dynamicLabels:     opts.DynamicLabels,
		dynamicLabelsFunc: labelValueFunc,
		valueFunc:         valueFunc,
		readyFunc:         opts.ReadyFunc,
	}

	return c
}

type DynamicLabelGaugeCollector struct {
	description       *prometheus.Desc
	dynamicLabels     []string
	dynamicLabelsFunc func() []string

	valueFunc func() float64
	readyFunc func() bool
}

func (c *DynamicLabelGaugeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.description
}

func (c *DynamicLabelGaugeCollector) Collect(ch chan<- prometheus.Metric) {
	if c.readyFunc != nil && !c.readyFunc() {
		return
	}
	ch <- prometheus.MustNewConstMetric(
		c.description,
		prometheus.GaugeValue,
		c.valueFunc(),
		c.dynamicLabelsFunc()...,
	)
}
