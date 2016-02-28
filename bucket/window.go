package bucket

import (
	"math"
	"sync"
)

const WindowSize = 1024

type Window struct {
	sync.RWMutex

	head    int
	size    int
	buckets []Interface

	cumSum float64

	last float64
}

func NewWindow(size int) *Window {
	w := &Window{
		size:    size,
		buckets: make([]Interface, size),
	}

	for i := 0; i < size; i++ {
		w.buckets[i] = &fixed{}
	}

	return w
}

func (w *Window) Push(b Interface) {
	w.Lock()
	defer w.Unlock()

	if w.buckets[w.head].Freq() > 0 {
		w.last = w.buckets[w.head].Last()
		w.cumSum += w.buckets[w.head].Sum()
	}

	w.buckets[w.head] = NewFixed(b)
	w.head = (w.head + 1) % w.size
}

func (w *Window) Count() []float64 {
	return w.mapFloat(func(b Interface) float64 {
		return b.Freq()
	})
}

func (w *Window) Sum() []float64 {
	return w.mapFloat(func(b Interface) float64 {
		return b.Sum()
	})
}

func (w *Window) CumSum() []float64 {
	sums := w.Sum()
	cumSum := w.cumSum
	for i, sum := range sums {
		sums[i] += cumSum
		cumSum += sum
	}
	return sums
}

func (w *Window) Unique() []float64 {
	return w.mapFloat(func(b Interface) float64 {
		return b.Unique()
	})
}

func (w *Window) Mean() []float64 {
	return w.mapFloat(func(b Interface) float64 {
		return b.Mean()
	})
}

func (w *Window) Median() []float64 {
	return w.mapFloat(func(b Interface) float64 {
		return b.Median()
	})
}

func (w *Window) P75() []float64 {
	return w.mapFloat(func(b Interface) float64 {
		return b.P75()
	})
}

func (w *Window) P95() []float64 {
	return w.mapFloat(func(b Interface) float64 {
		return b.P95()
	})
}

func (w *Window) P99() []float64 {
	return w.mapFloat(func(b Interface) float64 {
		return b.P99()
	})
}

func (w *Window) UniquePercent() []float64 {
	return w.mapFloat(func(b Interface) float64 {
		ct := b.Freq()
		if ct == 0 {
			return 0.0
		}
		return 100.0 * b.Unique() / ct
	})
}

func (w *Window) Last() []float64 {
	lasts := w.mapFloat(func(b Interface) float64 {
		if b.Freq() == 0 {
			return math.MaxFloat64
		}
		return b.Last()
	})

	if lasts[0] == math.MaxFloat64 {
		w.RLock()
		lasts[0] = w.last
		w.RUnlock()
	}

	for i := 1; i < w.size; i++ {
		if lasts[i] == math.MaxFloat64 {
			lasts[i] = lasts[i-1]
		}
	}

	return lasts
}

func (w *Window) mapFloat(f func(b Interface) float64) []float64 {
	w.RLock()
	defer w.RUnlock()

	vals := make([]float64, w.size)

	for i := 0; i < w.size; i++ {
		idx := (w.head + i) % w.size
		vals[i] = f(w.buckets[idx])
	}

	return vals
}
