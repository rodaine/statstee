package bucket

import (
	"math"
	"sync"
)

const WindowSize = 1024
const EmptyValue = math.MaxFloat64

type mapFunc func(idx int, b Interface, prev float64) float64

type Window struct {
	sync.RWMutex

	head    int
	size    int
	buckets []Interface

	cumSum float64

	last float64
}

type Averages struct {
	EWMA1, EWMA5, EWMA10 float64
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

func (w *Window) Index(offset int) int {
	return (w.head + offset) % w.size
}

func (w *Window) Push(b Interface) {
	w.Lock()
	defer w.Unlock()

	if w.buckets[w.head].Freq() > 0 {
		w.last = w.buckets[w.head].Last()
		w.cumSum += w.buckets[w.head].Sum()
	}

	w.buckets[w.head] = NewFixed(b)
	w.head = w.Index(1)
}

func (w *Window) Count() []float64        { return w.mapFloat(w._count) }
func (w *Window) CountAverages() Averages { return w.averages(w._count, false) }

func (w *Window) Unique() []float64        { return w.mapFloat(w._unique) }
func (w *Window) UniqueAverages() Averages { return w.averages(w._unique, false) }

func (w *Window) UniquePercent() []float64        { return w.mapFloat(w._uniquePercent) }
func (w *Window) UniquePercentAverages() Averages { return w.averages(w._uniquePercent, true) }

func (w *Window) Sum() []float64        { return w.mapFloat(w._sum) }
func (w *Window) SumAverages() Averages { return w.averages(w._sum, false) }

func (w *Window) Median() []float64        { return w.mapFloat(w._median) }
func (w *Window) MedianAverages() Averages { return w.averages(w._median, true) }

func (w *Window) P75() []float64        { return w.mapFloat(w._p75) }
func (w *Window) P75Averages() Averages { return w.averages(w._p75, true) }

func (w *Window) P95() []float64        { return w.mapFloat(w._p95) }
func (w *Window) P95Averages() Averages { return w.averages(w._p95, true) }

func (w *Window) Last() []float64        { return w.mapFloat(w._last) }
func (w *Window) LastAverages() Averages { return w.averages(w._last, true) }

func (w *Window) CumSum() []float64 {
	sums := w.Sum()
	cumSum := w.cumSum
	for i, sum := range sums {
		sums[i] += cumSum
		cumSum += sum
	}
	return sums
}

func (w *Window) mapFloat(f mapFunc) []float64 {
	w.RLock()
	defer w.RUnlock()

	vals := make([]float64, w.size)

	prev := EmptyValue
	for i := 0; i < w.size; i++ {
		vals[i] = f(i, w.buckets[w.Index(i)], prev)
		prev = vals[i]
	}

	return vals
}

func (w *Window) mapEWMA(f mapFunc, fill bool) []float64 {
	w.RLock()
	defer w.RUnlock()

	if w.size == 0 {
		return []float64{}
	}

	i := 0
	for ; i < w.size; i++ {
		if w.buckets[w.Index(i)].Freq() > 0 {
			break
		}
	}

	vals := make([]float64, w.size-i)
	prev := EmptyValue

	for j := 0; i < w.size; i, j = i+1, j+1 {
		b := w.buckets[w.Index(i)]

		if b.Freq() == 0 && fill && prev != EmptyValue {
			vals[j] = prev
			continue
		}

		vals[j] = f(i, b, prev)
		prev = vals[j]
	}

	return vals
}

// based off: https://en.wikipedia.org/wiki/Moving_average#Application_to_measuring_computer_performance
func (w *Window) ewma(data []float64, interval float64) float64 {
	ct := len(data)
	if ct == 0 {
		return 0.0
	}

	W := float64(ct) / 60.0 // How long is the data relevant (in minutes)
	a := math.Exp(-1.0 / (W * interval))

	y := data[0]
	for i := 1; i < ct; i++ {
		y = data[i] + a*(y-data[i])
	}

	return y
}

func (w *Window) averages(f mapFunc, fill bool) Averages {
	data := w.mapEWMA(f, fill)
	return Averages{
		EWMA1:  w.ewma(data, 1),
		EWMA5:  w.ewma(data, 5),
		EWMA10: w.ewma(data, 10),
	}
}

func (w *Window) _count(_ int, b Interface, _ float64) float64  { return b.Freq() }
func (w *Window) _unique(_ int, b Interface, _ float64) float64 { return b.Unique() }
func (w *Window) _sum(_ int, b Interface, _ float64) float64    { return b.Sum() }
func (w *Window) _median(_ int, b Interface, _ float64) float64 { return b.Median() }
func (w *Window) _p75(_ int, b Interface, _ float64) float64    { return b.P75() }
func (w *Window) _p95(_ int, b Interface, _ float64) float64    { return b.P95() }

func (w *Window) _uniquePercent(_ int, b Interface, _ float64) float64 {
	ct := b.Freq()
	if ct == 0 {
		return ct
	}
	return 100.0 * b.Unique() / ct
}

func (w *Window) _last(i int, b Interface, prev float64) float64 {
	if b.Freq() == 0 {
		if i == 0 {
			return w.last
		}
		return prev
	}
	return b.Last()
}
