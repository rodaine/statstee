package streams

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevice_IsLoopback(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Interface net.Interface
		Expected  bool
	}{
		{net.Interface{Flags: LoopbackFlags}, true},
		{net.Interface{Flags: net.FlagLoopback}, false},
		{net.Interface{Flags: net.FlagUp}, false},
		{net.Interface{}, false},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, isLoopback(test.Interface), "%+v", test)
	}
}

func TestDevice_ResolveDevice(t *testing.T) {
	t.Parallel()

	ifaces, _ := net.Interfaces()
	var loopback *net.Interface
	for _, iface := range ifaces {
		if isLoopback(iface) {
			loopback = &iface
			break
		}
	}
	if loopback == nil {
		t.Skip("no loopback network interface found")
	}

	i, err := resolveDevice("this is not a real device")
	assert.Error(t, err)

	i, err = resolveDevice(loopback.Name)
	assert.NoError(t, err)
	assert.Equal(t, i.Name, loopback.Name)

	i, err = resolveDevice(LoopbackAbbr)
	assert.NoError(t, err)
	assert.Equal(t, i.Name, loopback.Name)
}
