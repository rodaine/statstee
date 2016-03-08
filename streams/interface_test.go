package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterface_ResolveStream(t *testing.T) {
	t.Parallel()
	port := DefaultStatsDPort + 55 // arbitrary so as to not collide with other tests

	s, err := ResolveStream(ListenMode, "", port)
	assert.NoError(t, err)
	assert.IsType(t, &listener{}, s)
	s.(*listener).conn.Close()

	s, err = ResolveStream(CaptureMode, LoopbackAbbr, port)
	assert.NoError(t, err)
	assert.IsType(t, &sniffer{}, s)

	s, err = ResolveStream(DefaultMode, "this is an erroneous device", port)
	assert.Error(t, err)
	assert.Nil(t, s)

	s, err = ResolveStream(DefaultMode, LoopbackAbbr, port)
	assert.NoError(t, err)
	assert.IsType(t, &listener{}, s)
	defer s.(*listener).conn.Close()

	s, err = ResolveStream(DefaultMode, LoopbackAbbr, port)
	assert.NoError(t, err)
	assert.IsType(t, &sniffer{}, s)
}
