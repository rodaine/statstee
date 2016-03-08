package bucket

import (
	"testing"
	"time"

	"github.com/rodaine/statstee/datagram"
	"github.com/stretchr/testify/assert"
)

func TestMetricWindow(t *testing.T) {
	t.Parallel()

	w := NewMetricWindow(datagram.DummyMetric, 2, time.Second)
	assert.Zero(t, w.Last()[1])
	w.Add(1)
	<-time.NewTimer(1100 * time.Millisecond).C
	assert.Equal(t, 1.0, w.Last()[1])
}
