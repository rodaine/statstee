package streams

import (
	"fmt"
	"net"
)

const (
	// Placeholder to identify the first loopback device (usually lo or lo0, depending on platform)
	LoopbackAbbr = "_first_loopback_"

	// Device flags used to identify a loopback device
	LoopbackFlags = net.FlagLoopback | net.FlagUp
)

func resolveDevice(iface string) (i net.Interface, err error) {
	ifaces, _ := net.Interfaces()

	for _, i := range ifaces {
		if i.Name == iface {
			return i, nil
		}

		if iface == LoopbackAbbr && isLoopback(i) {
			return i, nil
		}
	}

	return i, fmt.Errorf("unknown interface device: %s", iface)
}

func isLoopback(iface net.Interface) bool {
	return iface.Flags&LoopbackFlags == LoopbackFlags
}
