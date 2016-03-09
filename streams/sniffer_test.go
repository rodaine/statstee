package streams

import (
	"os"
	"testing"

	"sync"

	"time"

	"github.com/google/gopacket/pcap"
	"github.com/rodaine/statstee/datagram"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var canPCAP = false

func TestMain(m *testing.M) {
	i, _ := resolveDevice(LoopbackAbbr)
	h, err := pcap.OpenLive(i.Name, 1, false, time.Second)

	canPCAP = err == nil
	if canPCAP {
		h.Close()
	}

	os.Exit(m.Run())
}

func TestSniffer_New(t *testing.T) {
	t.Parallel()
	port := DefaultStatsDPort + 111

	s, err := newSniffer("this is not a real device", port)
	assert.Error(t, err)
	assert.Nil(t, s)

	s, err = newSniffer(LoopbackAbbr, port)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	s.(*sniffer).handle.CleanUp()
}

func TestSniffer_Listen(t *testing.T) {
	t.Parallel()
	if !canPCAP {
		t.Skip("cannot use PCAP due to permissions. Run tests as sudo")
	}

	port := DefaultStatsDPort - 123

	s, err := newSniffer(LoopbackAbbr, port)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	var serr error
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		serr = s.Listen(ctx)
		wg.Done()
	}()

	m := datagram.Metric{
		Name:       "foo.bar",
		Type:       datagram.Counter,
		Value:      1,
		SampleRate: 1,
	}

	sender, _ := datagram.NewSender("localhost", port)
	sender.Send(m)

	b := <-s.Chan()
	assert.Equal(t, m.String(), string(b))

	cancel()
	wg.Wait()

	assert.NoError(t, serr)
}
