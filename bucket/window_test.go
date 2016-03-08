package bucket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWindow_Count(t *testing.T) {
	t.Parallel()

	w := NewWindow(3)

	tests := []struct {
		vals     []float64
		expected []float64
	}{
		{[]float64{}, []float64{0, 0, 0}},
		{[]float64{1}, []float64{0, 0, 1}},
		{[]float64{1, 2}, []float64{0, 1, 2}},
		{[]float64{1, 2, 3}, []float64{1, 2, 3}},
		{[]float64{}, []float64{2, 3, 0}},
	}

	for _, test := range tests {
		b := NewRaw()
		for _, v := range test.vals {
			b.Add(v)
		}
		w.Push(b)
		assert.EqualValues(t, test.expected, w.Count(), "%+v", test)
	}
}

func TestWindow_Sum(t *testing.T) {
	t.Parallel()

	w := NewWindow(3)

	tests := []struct {
		vals     []float64
		expected []float64
	}{
		{[]float64{}, []float64{0, 0, 0}},
		{[]float64{1}, []float64{0, 0, 1}},
		{[]float64{1, 2}, []float64{0, 1, 3}},
		{[]float64{1, 2, 3}, []float64{1, 3, 6}},
		{[]float64{}, []float64{3, 6, 0}},
	}

	for _, test := range tests {
		b := NewRaw()
		for _, v := range test.vals {
			b.Add(v)
		}
		w.Push(b)
		assert.EqualValues(t, test.expected, w.Sum(), "%+v", test)
	}
}

func TestWindow_CumSum(t *testing.T) {
	t.Parallel()

	w := NewWindow(3)

	tests := []struct {
		vals     []float64
		expected []float64
	}{
		{[]float64{}, []float64{0, 0, 0}},
		{[]float64{1}, []float64{0, 0, 1}},
		{[]float64{1, 2}, []float64{0, 1, 4}},
		{[]float64{1, 2, 3}, []float64{1, 4, 10}},
		{[]float64{1, 2}, []float64{4, 10, 13}},
		{[]float64{1}, []float64{10, 13, 14}},
		{[]float64{}, []float64{13, 14, 14}},
	}

	for _, test := range tests {
		b := NewRaw()
		for _, v := range test.vals {
			b.Add(v)
		}
		w.Push(b)
		assert.EqualValues(t, test.expected, w.CumSum(), "%+v", test)
	}
}

func TestWindow_Unique(t *testing.T) {
	t.Parallel()

	w := NewWindow(3)

	tests := []struct {
		vals     []float64
		expected []float64
	}{
		{[]float64{}, []float64{0, 0, 0}},
		{[]float64{1}, []float64{0, 0, 1}},
		{[]float64{1, 2}, []float64{0, 1, 2}},
		{[]float64{1, 2, 1}, []float64{1, 2, 2}},
		{[]float64{}, []float64{2, 2, 0}},
	}

	for _, test := range tests {
		b := NewRaw()
		for _, v := range test.vals {
			b.Add(v)
		}
		w.Push(b)
		assert.EqualValues(t, test.expected, w.Unique(), "%+v", test)
	}
}

func TestWindow_Median(t *testing.T) {
	t.Parallel()

	w := NewWindow(3)

	tests := []struct {
		vals     []float64
		expected []float64
	}{
		{[]float64{}, []float64{0, 0, 0}},
		{[]float64{1}, []float64{0, 0, 1}},
		{[]float64{1, 2}, []float64{0, 1, 1.5}},
		{[]float64{1, 2, 2, 10000}, []float64{1, 1.5, 2}},
		{[]float64{}, []float64{1.5, 2, 0}},
	}

	for _, test := range tests {
		b := NewRaw()
		for _, v := range test.vals {
			b.Add(v)
		}
		w.Push(b)
		assert.EqualValues(t, test.expected, w.Median(), "%+v", test)
	}
}

func TestWindow_P75(t *testing.T) {
	t.Parallel()

	w := NewWindow(3)

	tests := []struct {
		vals     []float64
		expected []float64
	}{
		{[]float64{}, []float64{0, 0, 0}},
		{[]float64{1}, []float64{0, 0, 1}},
		{[]float64{1, 2}, []float64{0, 1, 2}},
		{[]float64{1, 2, 3, 4}, []float64{1, 2, 3.75}},
		{[]float64{}, []float64{2, 3.75, 0}},
	}

	for _, test := range tests {
		b := NewRaw()
		for _, v := range test.vals {
			b.Add(v)
		}
		w.Push(b)
		assert.EqualValues(t, test.expected, w.P75(), "%+v", test)
	}
}

func TestWindow_P95(t *testing.T) {
	t.Parallel()

	w := NewWindow(3)

	tests := []struct {
		vals     []float64
		expected []float64
	}{
		{[]float64{}, []float64{0, 0, 0}},
		{[]float64{1}, []float64{0, 0, 1}},
		{[]float64{1, 2}, []float64{0, 1, 2}},
		{[]float64{1, 2, 3, 4}, []float64{1, 2, 4}},
		{[]float64{}, []float64{2, 4, 0}},
	}

	for _, test := range tests {
		b := NewRaw()
		for _, v := range test.vals {
			b.Add(v)
		}
		w.Push(b)
		assert.EqualValues(t, test.expected, w.P95(), "%+v", test)
	}
}

func TestWindow_UniquePercent(t *testing.T) {
	t.Parallel()

	w := NewWindow(3)

	tests := []struct {
		vals     []float64
		expected []float64
	}{
		{[]float64{}, []float64{0, 0, 0}},
		{[]float64{1}, []float64{0, 0, 100}},
		{[]float64{1, 2}, []float64{0, 100, 100}},
		{[]float64{1, 2, 1, 2}, []float64{100, 100, 50}},
		{[]float64{}, []float64{100, 50, 0}},
	}

	for _, test := range tests {
		b := NewRaw()
		for _, v := range test.vals {
			b.Add(v)
		}
		w.Push(b)
		assert.EqualValues(t, test.expected, w.UniquePercent(), "%+v", test)
	}
}

func TestWindow_Last(t *testing.T) {
	t.Parallel()

	w := NewWindow(3)

	tests := []struct {
		vals     []float64
		expected []float64
	}{
		{[]float64{}, []float64{0, 0, 0}},
		{[]float64{1}, []float64{0, 0, 1}},
		{[]float64{}, []float64{0, 1, 1}},
		{[]float64{2}, []float64{1, 1, 2}},
		{[]float64{}, []float64{1, 2, 2}},
		{[]float64{}, []float64{2, 2, 2}},
		{[]float64{}, []float64{2, 2, 2}},
		{[]float64{0}, []float64{2, 2, 0}},
	}

	for _, test := range tests {
		b := NewRaw()
		for _, v := range test.vals {
			b.Add(v)
		}
		w.Push(b)
		assert.EqualValues(t, test.expected, w.Last(), "%+v", test)
	}
}

func TestWindow_EWMA(t *testing.T) {
	t.Parallel()

	data := make([]float64, 600)
	for i := 0; i < len(data); i++ {
		data[i] = float64(i)
	}

	w := NewWindow(600)
	assert.Zero(t, w.ewma([]float64{}, 1))

	assert.InDelta(t, 589.5, w.ewma(data, 1), 0.1)
	assert.InDelta(t, 549.5, w.ewma(data, 5), 0.1)
	assert.InDelta(t, 499.7, w.ewma(data, 10), 0.1)
}

func TestWindow_MapEWMA(t *testing.T) {
	t.Parallel()

	w := NewWindow(0)
	assert.Empty(t, w.mapEWMA(w._sum, false))

	w = NewWindow(3)

	tests := []struct {
		vals     []float64
		fill     bool
		expected []float64
	}{
		{[]float64{}, false, []float64{}},
		{[]float64{1}, false, []float64{1}},
		{[]float64{2}, false, []float64{1, 2}},
		{[]float64{}, false, []float64{1, 2, 0}},
		{[]float64{3}, true, []float64{2, 2, 3}},
	}

	for _, test := range tests {
		b := NewRaw()
		for _, v := range test.vals {
			b.Add(v)
		}
		w.Push(b)
		assert.EqualValues(t, test.expected, w.mapEWMA(w._sum, test.fill))
	}
}

func TestWindow_Averages(t *testing.T) {
	t.Parallel()

	w := NewWindow(600)
	for i := 0; i < 600; i++ {
		b := NewRaw()
		b.Add(float64(i))
		w.Push(b)
	}

	avgs := w.averages(w._sum, false)
	assert.InDelta(t, 589.5, avgs.EWMA1, 0.1)
	assert.InDelta(t, 549.5, avgs.EWMA5, 0.1)
	assert.InDelta(t, 499.7, avgs.EWMA10, 0.1)
}

func TestWindow_AllAverages(t *testing.T) {
	t.Parallel()

	w := NewWindow(600)
	for i := 0; i < 600; i++ {
		b := NewRaw()
		b.Add(1)
		w.Push(b)
	}

	tests := []struct {
		f      func() Averages
		ewma1  float64
		ewma5  float64
		ewma10 float64
	}{
		{w.CountAverages, 1, 1, 1},
		{w.UniqueAverages, 1, 1, 1},
		{w.UniquePercentAverages, 100, 100, 100},
		{w.SumAverages, 1, 1, 1},
		{w.MedianAverages, 1, 1, 1},
		{w.P75Averages, 1, 1, 1},
		{w.P95Averages, 1, 1, 1},
		{w.LastAverages, 1, 1, 1},
	}

	for _, test := range tests {
		avgs := test.f()
		assert.Equal(t, test.ewma1, avgs.EWMA1, "%+v", test)
		assert.Equal(t, test.ewma5, avgs.EWMA5, "%+v", test)
		assert.Equal(t, test.ewma10, avgs.EWMA10, "%+v", test)
	}
}
