package bucket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixed(t *testing.T) {
	raw := NewRaw(1000)
	for i := 1; i <= 1000; i++ {
		raw.Add(float64(i))
	}

	fixed := NewFixed(raw)

	is := assert.New(t)
	is.Equal(raw.Sum(), fixed.Sum())
	is.Equal(raw.Freq(), fixed.Freq())
	is.Equal(raw.Unique(), fixed.Unique())
	is.Equal(raw.Last(), fixed.Last())
	is.Equal(raw.Mean(), fixed.Mean())
	is.Equal(raw.Min(), fixed.Min())
	is.Equal(raw.Max(), fixed.Max())
	is.Equal(raw.Median(), fixed.Median())
	is.Equal(raw.P75(), fixed.P75())
	is.Equal(raw.P95(), fixed.P95())
	is.Equal(raw.P99(), fixed.P99())

	fixed.Add(123)
	is.Equal(raw.Freq(), fixed.Freq())

	fixed.Reset()

	is.Zero(fixed.Sum())
	is.Zero(fixed.Freq())
	is.Zero(fixed.Last())
}
