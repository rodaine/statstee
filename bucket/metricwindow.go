package bucket

import (
	"sync"
	"time"

	"github.com/rodaine/statstee/datagram"
)

var DummyWindow = NewMetricWindow(datagram.DummyMetric, 1, time.Hour)

type MetricWindow struct {
	sync.Mutex
	*Window

	Metric datagram.Metric
	curr   Interface
}

func NewMetricWindow(m datagram.Metric, size int, d time.Duration) *MetricWindow {
	w := &MetricWindow{
		Metric: m,
		Window: NewWindow(size),
		curr:   NewRaw(),
	}

	go w.tick(d)

	return w
}

func (w *MetricWindow) tick(d time.Duration) {
	for range time.NewTicker(d).C {
		w.Lock()

		w.Push(w.curr)
		w.curr.Reset()

		w.Unlock()
	}
}

func (w *MetricWindow) Add(v float64) {
	w.Lock()
	defer w.Unlock()
	w.curr.Add(v)
}
