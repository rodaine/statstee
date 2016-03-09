package streams

import (
	"fmt"
	"net"

	"log"

	"time"

	"golang.org/x/net/context"
)

type listener struct {
	conn *net.UDPConn
	c    chan []byte
}

func newListener(port int) (Interface, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("unable to resolve UDP address: %v", err)
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Printf("unable to open UDP connection: %v", err)
		return nil, err
	}

	return &listener{
		conn: conn,
		c:    make(chan []byte, channelBuffer),
	}, nil
}

func (l *listener) Listen(ctx context.Context) error {
	defer l.conn.Close()
	defer close(l.c)

	buffer := make([]byte, maxDatagramSize)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			//noop
		}

		l.conn.SetReadDeadline(time.Now().Add(time.Second))
		n, _, err := l.conn.ReadFromUDP(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok {
				if netErr.Timeout() || netErr.Temporary() {
					continue
				}
			}
			log.Printf("unable to read bytes from connection: %v", err)
			return err
		}

		raw := make([]byte, n)
		copy(raw, buffer[:n])

		l.c <- raw
	}
}

func (l *listener) Chan() <-chan []byte {
	return l.c
}

var _ Interface = &listener{}
