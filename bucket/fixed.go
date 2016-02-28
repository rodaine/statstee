package bucket

type fixed struct {
	sum                   float64
	freq, uniq            float64
	last                  float64
	mean                  float64
	min, max              float64
	median, p75, p95, p99 float64
}

func NewFixed(b Interface) Interface {
	return &fixed{
		sum:    b.Sum(),
		freq:   b.Freq(),
		uniq:   b.Unique(),
		last:   b.Last(),
		mean:   b.Mean(),
		min:    b.Min(),
		max:    b.Max(),
		median: b.Median(),
		p75:    b.P75(),
		p95:    b.P95(),
		p99:    b.P99(),
	}
}

func (b *fixed) Add(m float64) {
	// noop
}

func (b *fixed) Sum() float64 {
	return b.sum
}

func (b *fixed) Freq() float64 {
	return b.freq
}

func (b *fixed) Unique() float64 {
	return b.uniq
}

func (b *fixed) Last() float64 {
	return b.last
}

func (b *fixed) Mean() float64 {
	return b.mean
}

func (b *fixed) Min() float64 {
	return b.min
}

func (b *fixed) Max() float64 {
	return b.max
}

func (b *fixed) Median() float64 {
	return b.median
}

func (b *fixed) P75() float64 {
	return b.p75
}

func (b *fixed) P95() float64 {
	return b.p95
}

func (b *fixed) P99() float64 {
	return b.p99
}

func (b *fixed) Reset() {
	b.sum = 0
	b.freq = 0
	b.last = 0
	b.mean = 0
	b.min = 0
	b.max = 0
	b.min = 0
	b.max = 0
	b.median = 0
	b.p75 = 0
	b.p95 = 0
	b.p99 = 0
}

var _ Interface = &fixed{}
