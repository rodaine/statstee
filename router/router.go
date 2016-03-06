package router

import (
	"sync"

	"time"

	"sort"

	"github.com/rodaine/statstee/bucket"
	"github.com/rodaine/statstee/datagram"
)

type Router struct {
	sync.RWMutex

	c <-chan datagram.Metric

	metrics       []*bucket.MetricWindow
	metricsLookup map[string]int

	selected string
}

func (r *Router) Len() int           { return len(r.metrics) }
func (r *Router) Less(i, j int) bool { return r.metrics[i].Metric.Name < r.metrics[j].Metric.Name }
func (r *Router) Swap(i, j int) {
	a, b := r.metrics[i].Metric.Name, r.metrics[j].Metric.Name
	r.metrics[i], r.metrics[j] = r.metrics[j], r.metrics[i]
	r.metricsLookup[a], r.metricsLookup[b] = j, i
}

func New(c <-chan datagram.Metric) *Router {
	r := &Router{
		c:             c,
		metrics:       []*bucket.MetricWindow{},
		metricsLookup: map[string]int{},
	}
	return r
}

func (r *Router) Listen() {
	for m := range r.c {
		w := r.addOrGet(m)
		w.Add(m.Value)
	}
}

func (r *Router) addOrGet(m datagram.Metric) *bucket.MetricWindow {
	r.RLock()
	if idx, found := r.metricsLookup[m.Name]; found {
		r.RUnlock()
		return r.metrics[idx]
	}
	r.RUnlock()

	return r.Add(m)
}

func (r *Router) Add(m datagram.Metric) *bucket.MetricWindow {
	r.Lock()
	defer r.Unlock()

	w := bucket.NewMetricWindow(m, bucket.WindowSize, time.Second)

	idx := len(r.metrics)
	r.metrics = append(r.metrics, w)
	r.metricsLookup[m.Name] = idx
	sort.Sort(r)

	if r.selected == "" {
		r.selected = m.Name
	}

	return w
}

func (r *Router) Selected() string {
	r.RLock()
	defer r.RUnlock()
	return r.selected
}

func (r *Router) Metrics() []datagram.Metric {
	r.RLock()
	defer r.RUnlock()

	out := make([]datagram.Metric, len(r.metrics))
	for i, m := range r.metrics {
		out[i] = m.Metric
	}

	return out
}

func (r *Router) SelectedMetric() *bucket.MetricWindow {
	r.RLock()
	defer r.RUnlock()

	if r.selected == "" {
		return bucket.DummyWindow
	}

	idx, found := r.metricsLookup[r.selected]
	if !found {
		return bucket.DummyWindow
	}

	return r.metrics[idx]
}

func (r *Router) Previous() {
	r.Lock()
	defer r.Unlock()

	if len(r.metrics) == 0 {
		r.selected = ""
		return
	}

	idx, found := r.metricsLookup[r.selected]
	if r.selected == "" || !found {
		r.selected = r.metrics[0].Metric.Name
		return
	}

	if idx == 0 {
		return
	}

	r.selected = r.metrics[idx-1].Metric.Name
}

func (r *Router) Next() {
	r.Lock()
	defer r.Unlock()

	if len(r.metrics) == 0 {
		r.selected = ""
		return
	}

	idx, found := r.metricsLookup[r.selected]
	if r.selected == "" || !found {
		r.selected = r.metrics[0].Metric.Name
		return
	}

	if idx == len(r.metrics)-1 {
		return
	}

	r.selected = r.metrics[idx+1].Metric.Name
}
