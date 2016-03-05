package streams

import (
	"fmt"
	"net"

	"log"

	"golang.org/x/net/context"
)

type listener struct {
	conn *net.UDPConn
	c    chan []byte
}

func NewListener(port int) (Interface, error) {
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
		c:    make(chan []byte, 1000),
	}, nil
}

func (l *listener) Listen(ctx context.Context) error {
	defer l.conn.Close()
	defer close(l.c)

	buffer := make([]byte, MaxDatagramSize)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			//noop
		}

		n, _, err := l.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("unable to read bytes from connection: %v", err)
			continue
		}

		raw := make([]byte, n)
		copy(raw, buffer[:n])

		l.c <- raw
	}
}

func (l *listener) Chan() <-chan []byte {
	return l.c
}
