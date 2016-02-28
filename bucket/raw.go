package bucket

import (
	"math"
	"sort"
)

type raw struct {
	values []float64
	last   float64
}

func NewRaw() Interface {
	return &raw{values: []float64{}}
}

func (b *raw) Add(m float64) {
	b.values = append(b.values, m)
	b.last = m
}

func (b *raw) Sum() float64 {
	sum := 0.0
	for _, val := range b.values {
		sum += val
	}

	return sum
}

func (b *raw) Freq() float64 {
	return float64(len(b.values))
}

func (b *raw) Unique() float64 {
	freq := len(b.values)
	if freq == 0 {
		return 0
	}

	b.sort()

	ct := 1
	for i := 1; i < freq; i++ {
		if b.values[i] != b.values[i-1] {
			ct++
		}
	}

	return float64(ct)
}

func (b *raw) Last() float64 {
	return b.last
}

func (b *raw) Mean() float64 {
	ct := b.Freq()
	if ct == 0 {
		return 0.0
	}

	return b.Sum() / float64(ct)
}

func (b *raw) Min() float64 {
	if b.Freq() == 0 {
		return 0.0
	}

	b.sort()
	return b.values[0]
}

func (b *raw) Max() float64 {

	ct := len(b.values)
	if ct == 0 {
		return 0.0
	}

	b.sort()
	return b.values[ct-1]
}

func (b *raw) Median() float64 {

	return b.percentile(0.5)
}

func (b *raw) P75() float64 {

	return b.percentile(0.75)
}

func (b *raw) P95() float64 {

	return b.percentile(0.95)
}

func (b *raw) P99() float64 {
	return b.percentile(0.99)
}

func (b *raw) Reset() {
	b.values = make([]float64, 0, cap(b.values))
	b.last = 0.0
}

func (b *raw) sort() {
	if !sort.Float64sAreSorted(b.values) {
		sort.Float64s(b.values)
	}
}

func (b *raw) percentile(p float64) float64 {
	if p <= 0.0 || p > 1.0 {
		panic("percentile out of range")
	}

	ct := len(b.values)
	if ct == 0 {
		return 0.0
	}

	b.sort()

	n := float64(ct)
	if p >= n/(n+1) {
		return b.values[ct-1]
	}

	idx := p * (n + 1)
	i := int(idx)
	r := math.Mod(idx, 1)
	return b.values[i-1] + r*(b.values[i]-b.values[i-1])
}

var _ Interface = &raw{}
