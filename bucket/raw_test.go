package bucket

import (
	"testing"

	"math"

	"github.com/stretchr/testify/assert"
)

func TestRaw_Sum(t *testing.T) {
	t.Parallel()

	sum := 0.0
	b := NewRaw()
	assert.Zero(t, b.Sum())

	for _, v := range []float64{1, -2, 3, -4, 5} {
		sum += v
		b.Add(v)
	}

	assert.Equal(t, sum, b.Sum())
}

func TestRaw_Freq(t *testing.T) {
	t.Parallel()

	ct := 0
	b := NewRaw()
	assert.Zero(t, b.Freq())

	for _, v := range []float64{1, 2, 3, 4, 5} {
		ct++
		b.Add(v)
	}

	assert.Equal(t, float64(ct), b.Freq())
}

func TestRaw_Unique(t *testing.T) {
	t.Parallel()

	b := NewRaw()
	assert.Zero(t, b.Unique())

	for _, v := range []float64{1, 1, 2, 2, 3} {
		b.Add(v)
	}

	assert.Equal(t, float64(3), b.Unique())
}

func TestRaw_Last(t *testing.T) {
	t.Parallel()

	var last float64
	b := NewRaw()
	for _, last = range []float64{1, 3, 5, 7, 9} {
		b.Add(last)
	}

	assert.Equal(t, last, b.Last())
}

func TestRaw_Mean(t *testing.T) {
	t.Parallel()

	b := NewRaw()
	assert.Zero(t, b.Mean())

	ct := 0
	sum := 0.0

	for _, v := range []float64{1, 2, 3, 4, 5} {
		ct++
		sum += v
		b.Add(v)
	}

	assert.Equal(t, sum/float64(ct), b.Mean())
}

func TestRaw_Min(t *testing.T) {
	t.Parallel()

	b := NewRaw()
	assert.Zero(t, b.Min())

	min := math.MaxFloat64
	for _, v := range []float64{5, 1, 4, 3, 2} {
		min = math.Min(min, v)
		b.Add(v)
	}

	assert.Equal(t, min, b.Min())
}

func TestRaw_Max(t *testing.T) {
	t.Parallel()

	b := NewRaw()
	assert.Zero(t, b.Max())

	max := 0.0
	for _, v := range []float64{5, 1, 4, 3, 2} {
		max = math.Max(max, v)
		b.Add(v)
	}

	assert.Equal(t, max, b.Max())
}

func TestRaw_Median(t *testing.T) {
	t.Parallel()

	tests := []struct {
		vals     []float64
		expected float64
	}{
		{[]float64{}, 0},
		{[]float64{1, 2}, 1.5},
		{[]float64{1, 1, 1}, 1},
		{[]float64{3, 2, 1}, 2},
		{[]float64{4, 2, 3, 1}, 2.5},
	}

	for _, test := range tests {
		b := NewRaw()
		for _, v := range test.vals {
			b.Add(v)
		}
		assert.Equal(t, test.expected, b.Median(), "%+v", test)
	}
}

func TestRaw_P75(t *testing.T) {
	t.Parallel()

	tests := []struct {
		vals     []float64
		expected float64
	}{
		{[]float64{}, 0},
		{[]float64{1, 1, 1, 1, 1}, 1},
		{[]float64{3, 2, 1}, 3},
		{[]float64{4, 2, 3, 1}, 3.75},
		{[]float64{1, 2, 3, 4, 5}, 4.5},
	}

	for _, test := range tests {
		b := NewRaw()
		for _, v := range test.vals {
			b.Add(v)
		}
		assert.Equal(t, test.expected, b.P75(), "%+v", test)
	}
}

func TestRaw_P95(t *testing.T) {
	t.Parallel()

	vals := make([]float64, 100)
	for i := 1; i <= 100; i++ {
		vals[i-1] = float64(i)
	}

	tests := []struct {
		vals     []float64
		expected float64
	}{
		{[]float64{}, 0},
		{[]float64{3, 2, 1}, 3},
		{vals, 96},
	}

	for _, test := range tests {
		b := NewRaw()
		for _, v := range test.vals {
			b.Add(v)
		}
		assert.InDelta(t, test.expected, b.P95(), 0.1, "%+v", test)
	}
}

func TestRaw_P99(t *testing.T) {
	t.Parallel()

	vals := make([]float64, 1000)
	for i := 1; i <= 1000; i++ {
		vals[i-1] = float64(i)
	}

	tests := []struct {
		vals     []float64
		expected float64
	}{
		{[]float64{}, 0},
		{[]float64{3, 2, 1}, 3},
		{vals, 991},
	}

	for _, test := range tests {
		b := NewRaw()
		for _, v := range test.vals {
			b.Add(v)
		}
		assert.InDelta(t, test.expected, b.P99(), 0.1, "%+v", test)
	}
}

func TestRaw_Reset(t *testing.T) {
	t.Parallel()

	b := NewRaw()
	b.Add(123)

	is := assert.New(t)
	is.NotZero(b.Last())
	is.NotZero(b.Sum())
	is.NotZero(b.Freq())

	b.Reset()
	is.Zero(b.Last())
	is.Zero(b.Sum())
	is.Zero(b.Freq())
}

func TestRaw_PercentilePanic(t *testing.T) {
	t.Parallel()

	for _, p := range []float64{-1, 0, 1.1} {
		assert.Panics(t, func() {
			b, ok := NewRaw().(*raw)
			if !ok {
				assert.FailNow(t, "not a *raw")
				return
			}
			b.Add(0)
			b.percentile(p)
		})
	}
}
