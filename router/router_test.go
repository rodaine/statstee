package router

import (
	"testing"

	"time"

	"github.com/rodaine/statstee/bucket"
	"github.com/rodaine/statstee/datagram"
	"github.com/stretchr/testify/assert"
)

func TestRouter_Listen(t *testing.T) {
	t.Parallel()

	c := make(chan datagram.Metric, 3)
	c <- datagram.Metric{
		Name:       "foo.bar",
		Type:       datagram.Counter,
		Value:      1,
		SampleRate: 1,
	}
	c <- datagram.Metric{
		Name:       "fizz.buzz",
		Type:       datagram.Histogram,
		Value:      1,
		SampleRate: 1,
	}
	c <- datagram.Metric{
		Name:       "fizz.buzz",
		Type:       datagram.Histogram,
		Value:      1,
		SampleRate: 1,
	}
	close(c)

	r := New(c)
	r.Listen()
	<-time.NewTimer(100 * time.Millisecond).C

	m := r.Metrics()
	assert.Len(t, m, 2)

	assert.Equal(t, "foo.bar", r.Selected(), "first metric added should be selected")
	assert.Equal(t, "fizz.buzz", m[0].Name, "alphabetized metrics")
}

func TestRouter_SelectedMetric(t *testing.T) {
	t.Parallel()

	c := make(chan datagram.Metric)
	defer close(c)

	r := New(c)
	go r.Listen()

	assert.True(t, bucket.DummyWindow == r.SelectedMetric())
	r.selected = "foo"
	assert.True(t, bucket.DummyWindow == r.SelectedMetric())
	r.selected = ""

	c <- datagram.Metric{
		Name:       "foo.bar",
		Type:       datagram.Counter,
		Value:      1,
		SampleRate: 1,
	}

	assert.False(t, bucket.DummyWindow == r.SelectedMetric())
	assert.NotNil(t, r.SelectedMetric())
}

func TestRouter_PreviousNext(t *testing.T) {
	t.Parallel()

	c := make(chan datagram.Metric, 2)
	r := New(c)

	r.Previous()
	assert.Empty(t, r.Selected())
	r.Next()
	assert.Empty(t, r.Selected())

	c <- datagram.Metric{
		Name:       "foo.bar",
		Type:       datagram.Counter,
		Value:      1,
		SampleRate: 1,
	}
	c <- datagram.Metric{
		Name:       "fizz.buzz",
		Type:       datagram.Set,
		Value:      1,
		SampleRate: 1,
	}
	close(c)
	r.Listen()


	assert.Equal(t, "foo.bar", r.Selected(), "first metric added should be selected")
	r.selected = "not.found"
	r.Previous()
	assert.Equal(t, "fizz.buzz", r.Selected(), "not found metric should default to first")
	r.selected = ""
	r.Next()
	assert.Equal(t, "fizz.buzz", r.Selected(), "not found metric should default to first")

	r.Previous()
	assert.Equal(t, "fizz.buzz", r.Selected(), "if on first item, don't change on previous")
	r.Next()
	assert.Equal(t, "foo.bar", r.Selected())
	r.Next()
	assert.Equal(t, "foo.bar", r.Selected(), "if on last item, don't change on next")
	r.Previous()
	assert.Equal(t, "fizz.buzz", r.Selected())
}
