package streams

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"golang.org/x/net/context"
)

const (
	udpPortFilter = "udp port %d"
)

type sniffer struct {
	handle *pcap.InactiveHandle
	port   int
	c      chan []byte
}

func newSniffer(device string, port int) (Interface, error) {
	iface, err := resolveDevice(device)
	if err != nil {
		log.Printf("unable to resolve device: %v", err)
		return nil, err
	}

	handle, err := pcap.NewInactiveHandle(iface.Name)
	if err != nil {
		log.Printf("unable to create handle: %v", err)
		return nil, err
	}

	handle.SetSnapLen(maxDatagramSize)
	handle.SetImmediateMode(true)
	handle.SetPromisc(true)
	handle.SetTimeout(pcap.BlockForever)
	handle.SetRFMon(true)

	s := &sniffer{
		handle: handle,
		port:   port,
		c:      make(chan []byte, channelBuffer),
	}

	return s, nil
}

func (s *sniffer) Listen(ctx context.Context) error {
	defer s.handle.CleanUp()
	defer close(s.c)

	h, err := s.handle.Activate()
	if err != nil {
		err = fmt.Errorf("unable to activate pcap handle: %v", err)
		log.Println(err)
		return err
	}
	defer h.Close()

	if err = h.SetBPFFilter(fmt.Sprintf(udpPortFilter, s.port)); err != nil {
		err = fmt.Errorf("unable to apply BPF filter: %v", err)
		log.Println(err)
		return err
	}

	packetSource := gopacket.NewPacketSource(h, h.LinkType())
	for {
		select {
		case <-ctx.Done():
			return nil
		case packet := <-packetSource.Packets():
			raw := packet.ApplicationLayer().Payload()
			s.c <- raw
		}
	}
}

func (s *sniffer) Chan() <-chan []byte {
	return s.c
}

var _ Interface = &sniffer{}
