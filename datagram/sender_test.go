package datagram

import (
	"net"
	"testing"

	"sync"

	"github.com/stretchr/testify/assert"
)

func TestSender_BadHost(t *testing.T) {
	t.Parallel()

	s, err := NewSender("this is not a real host", -1)
	assert.Nil(t, s)
	assert.Error(t, err)
}

func TestSender_Send(t *testing.T) {
	t.Parallel()

	s, err := NewSender("localhost", 8182)
	assert.NoError(t, err)
	defer s.conn.Close()

	m := Metric{
		Name:       "test.sender",
		Type:       Counter,
		Value:      1,
		SampleRate: 1,
	}
	b := make([]byte, len(m.String()))

	addr, _ := net.ResolveUDPAddr("udp", "localhost:8182")
	conn, _ := net.ListenUDP("udp", addr)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		conn.ReadFromUDP(b)
		wg.Done()
		conn.Close()
	}()

	err = s.Send(m)
	assert.NoError(t, err)

	wg.Wait()

	assert.Equal(t, m.String(), string(b))
}
