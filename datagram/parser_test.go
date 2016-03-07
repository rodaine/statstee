package datagram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Parallel()

	data := make(chan []byte, 3)
	data <- []byte("test.parser:1|c")
	data <- []byte("totally malformed metric")
	data <- []byte("test.parser:2|c")
	close(data)

	first := Metric{
		Name:       "test.parser",
		Type:       Counter,
		Value:      1,
		SampleRate: 1,
	}

	second := Metric{
		Name:       "test.parser",
		Type:       Counter,
		Value:      2,
		SampleRate: 1,
	}

	p := NewParser()
	p.Parse(data)

	c := p.Chan()
	assert.EqualValues(t, first, <-c)
	assert.EqualValues(t, second, <-c)

	_, more := <-c
	assert.False(t, more)
}
