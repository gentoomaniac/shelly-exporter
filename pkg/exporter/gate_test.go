package exporter

import (
	"sync/atomic"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestGatedCollectorWithholdsMetricsWhileOffline(t *testing.T) {
	online := &atomic.Bool{}
	inner := prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "test_value", Help: "test"}, func() float64 { return 42 })

	reg := prometheus.NewPedanticRegistry()
	if err := reg.Register(gatedCollector{Collector: inner, online: online}); err != nil {
		t.Fatalf("register: %v", err)
	}

	count := func() int {
		n, err := testutil.GatherAndCount(reg, "test_value")
		if err != nil {
			t.Fatalf("gather: %v", err)
		}
		return n
	}

	if n := count(); n != 0 {
		t.Errorf("before first successful refresh: got %d series, want 0", n)
	}

	online.Store(true)
	if n := count(); n != 1 {
		t.Errorf("online: got %d series, want 1", n)
	}

	online.Store(false)
	if n := count(); n != 0 {
		t.Errorf("offline: got %d series, want 0", n)
	}
}
