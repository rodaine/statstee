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
	needsUpdate bool

	c <-chan datagram.Metric

	windows      []*bucket.MetricWindow
	metricLookup map[string]int

	selected string
}

func (r Router) Len() int           { return len(r.windows) }
func (r Router) Less(i, j int) bool { return r.windows[i].Metric.Name < r.windows[j].Metric.Name }
func (r Router) Swap(i, j int) {
	a, b := r.windows[i].Metric.Name, r.windows[j].Metric.Name
	r.windows[i], r.windows[j] = r.windows[j], r.windows[i]
	r.metricLookup[a], r.metricLookup[b] = j, i
}

func New(c <-chan datagram.Metric) *Router {
	r := &Router{
		needsUpdate:  true,
		c:            c,
		windows:      []*bucket.MetricWindow{},
		metricLookup: map[string]int{},
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
	if idx, found := r.metricLookup[m.Name]; found {
		r.RUnlock()
		return r.windows[idx]
	}
	r.RUnlock()

	return r.Add(m)
}

func (r *Router) Add(m datagram.Metric) *bucket.MetricWindow {
	r.Lock()

	w := bucket.NewMetricWindow(m, bucket.WindowSize, time.Second)

	idx := len(r.windows)
	r.windows = append(r.windows, w)
	r.metricLookup[m.Name] = idx
	sort.Sort(r)

	if r.selected == "" {
		r.selected = m.Name
	}

	r.needsUpdate = true
	r.Unlock()
	return w
}

func (r *Router) Selected() string {
	r.RLock()
	s := r.selected
	r.RUnlock()

	return s
}

func (r *Router) NeedsUpdate() bool {
	r.RLock()
	u := r.needsUpdate
	r.RUnlock()

	if u {
		r.Lock()
		r.needsUpdate = false
		r.Unlock()
	}

	return u
}

func (r *Router) Metrics() []datagram.Metric {
	r.RLock()

	out := make([]datagram.Metric, len(r.windows))
	for i, m := range r.windows {
		out[i] = m.Metric
	}

	defer r.RUnlock()
	return out
}

func (r *Router) SelectedMetric() *bucket.MetricWindow {
	r.RLock()
	defer r.RUnlock()

	if r.selected == "" {
		return bucket.DummyWindow
	}

	idx, found := r.metricLookup[r.selected]
	if !found {
		return bucket.DummyWindow
	}

	return r.windows[idx]
}

func (r *Router) Previous() {
	r.Lock()
	defer r.Unlock()

	if len(r.windows) == 0 {
		r.selected = ""
		return
	}

	idx, found := r.metricLookup[r.selected]
	if r.selected == "" || !found {
		r.selected = r.windows[0].Metric.Name
		return
	}

	if idx == 0 {
		return
	}

	r.selected = r.windows[idx-1].Metric.Name
	r.needsUpdate = true
}

func (r *Router) Next() {
	r.Lock()
	defer r.Unlock()

	if len(r.windows) == 0 {
		r.selected = ""
		return
	}

	idx, found := r.metricLookup[r.selected]
	if r.selected == "" || !found {
		r.selected = r.windows[0].Metric.Name
		return
	}

	if idx == len(r.windows)-1 {
		return
	}

	r.selected = r.windows[idx+1].Metric.Name
	r.needsUpdate = true
}
