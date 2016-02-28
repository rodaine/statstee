package datagram

import (
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

const (
	DatagramSize  = 1024
	LoopbackAbbr  = "lo"
	LoopbackFlags = net.FlagLoopback | net.FlagUp
)

func Stream(iface string, port int, c chan<- Metric) error {
	iface, err := resolveDevice(iface)
	if err != nil {
		return err
	}

	handle, err := pcap.OpenLive(iface, DatagramSize, false, pcap.BlockForever)
	if err != nil {
		return err
	}
	defer handle.Close()

	filter := fmt.Sprintf("udp and port %d", port)
	if err := handle.SetBPFFilter(filter); err != nil {
		return fmt.Errorf("coult not set packet filter: %v", err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		raw := packet.ApplicationLayer().Payload()
		if m, err := ParseMetric(string(raw)); err == nil {
			c <- m
		}
	}

	return handle.Error()
}

func resolveDevice(iface string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return iface, err
	}

	for _, i := range ifaces {
		if i.Name == iface {
			return iface, nil
		}
		if iface == LoopbackAbbr && i.Flags&LoopbackFlags == LoopbackFlags {
			return i.Name, nil
		}
	}

	return iface, fmt.Errorf("unknown interface device: %s", iface)
}
