package datagram

import (
	"fmt"
	"net"
)

type Sender struct {
	conn *net.UDPConn
}

func NewSender(host string, port int) (*Sender, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &Sender{conn}, nil
}

func (s *Sender) Send(m Metric) error {
	_, err := s.conn.Write([]byte(m.String()))
	return err
}
