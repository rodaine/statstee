package streams

import (
	"testing"

	"sync"

	"time"

	"github.com/rodaine/statstee/datagram"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestListener_New(t *testing.T) {
	t.Parallel()
	port := DefaultStatsDPort - 42

	l, err := newListener(-1)
	assert.Error(t, err)
	assert.Nil(t, l)

	l, err = newListener(port)
	assert.NoError(t, err)
	assert.IsType(t, &listener{}, l)
	defer l.(*listener).conn.Close()

	l, err = newListener(port)
	assert.Error(t, err)
	assert.Nil(t, l)
}

func TestListener_Listen(t *testing.T) {
	t.Parallel()
	port := DefaultStatsDPort + 12

	ctx, cancel := context.WithCancel(context.Background())

	l, err := newListener(port)
	assert.NoError(t, err)
	assert.NotNil(t, l)

	var lerr error
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		lerr = l.Listen(ctx)
		wg.Done()
	}()

	m := datagram.Metric{
		Name:       "foo.bar",
		Type:       datagram.Counter,
		Value:      1,
		SampleRate: 1,
	}

	s, _ := datagram.NewSender("localhost", port)
	s.Send(m)

	b := <-l.Chan()
	cancel()

	assert.Equal(t, m.String(), string(b))
	wg.Wait()
	assert.NoError(t, lerr)
}

func TestListener_Listen_Error(t *testing.T) {
	t.Parallel()
	port := DefaultStatsDPort + 432

	l, err := newListener(port)
	assert.NoError(t, err)
	assert.NotNil(t, l)
	l.(*listener).conn.Close()

	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	err = l.Listen(ctx)
	assert.Error(t, err)
}
