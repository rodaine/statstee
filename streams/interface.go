package streams

import "golang.org/x/net/context"

const (
	// DefaultStatsDPort is the default port the StatsD daemon listens on
	DefaultStatsDPort = 8125

	// DatagramSize is the maximum size of a UDP datagram
	MaxDatagramSize = 8192
)

type StreamMode uint8

const (
	DefaultMode StreamMode = iota
	ListenMode
	CaptureMode
)

// Interface describes a stream source for stats metrics
type Interface interface {
	// Listen begins reading metric datagrams off the network, sending the raw bytes to the data channel. This method
	// blocks until ctx is Done or an internal error arises.
	Listen(ctx context.Context) error

	// Chan returns the channel through which the raw datagrams will be returned. If the channel is closed, this stream
	// is no longer valid and a new one will need to be created.
	Chan() <-chan []byte
}

func ResolveStream(mode StreamMode, device string, port int) (stream Interface, err error) {
	switch mode {
	case ListenMode:
		return NewListener(port)
	case CaptureMode:
		return NewSniffer(device, port)
	case DefaultMode:
		fallthrough
	default:
		if stream, err = NewListener(port); err != nil {
			return NewSniffer(device, port)
		}
		return
	}
}
