package bucket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWindow_Count(t *testing.T) {
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

func TestWindow_Unique(t *testing.T) {
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
