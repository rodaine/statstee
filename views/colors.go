package views

import (
	"fmt"

	"github.com/gizak/termui"
	"github.com/rodaine/statstee/datagram"
)

const (
	markdownColorFormat = "fg-%s,bg-%s"
)

var (
	datagramColorNames = map[datagram.MetricType]string{
		datagram.Histogram: "green",
		datagram.Timer:     "green",
		datagram.Counter:   "blue",
		datagram.Gauge:     "yellow",
		datagram.Set:       "magenta",
		datagram.Unknown:   "red",
	}

	datagramColors map[datagram.MetricType]termui.Attribute
)

func init() {
	datagramColors = make(map[datagram.MetricType]termui.Attribute, len(datagramColorNames))
	for t, name := range datagramColorNames {
		datagramColors[t] = termui.StringToAttribute(name)
	}
}

func markdown(fg, bg string) string {
	return fmt.Sprintf(markdownColorFormat, fg, bg)
}
