package views

import (
	"github.com/gizak/termui"
	"github.com/rodaine/statstee/bucket"
	"github.com/rodaine/statstee/datagram"
)

const (
	headerBorder     = termui.ColorWhite
	headerLabelColor = termui.ColorBlack | termui.AttrBold
)

type plotSet struct {
	metric datagram.Metric
	plots  []*plot
	header *termui.Par
	row    *termui.Row
}

func newSet(w *bucket.MetricWindow, prev *termui.Row) *plotSet {
	s := &plotSet{metric: w.Metric}
	s.header = s.setHeader()

	color := datagramColors[w.Metric.Type]

	switch w.Metric.Type {
	case datagram.Gauge:
		s.plots = []*plot{
			newPlot("Gauge Value", color, w.Last, true),
		}
	case datagram.Counter:
		s.plots = []*plot{
			newPlot("Count", color, w.Sum, true),
			newPlot("Cumulative Count", color, w.CumSum, false),
		}
	case datagram.Histogram, datagram.Timer:
		s.plots = []*plot{
			newPlot("Count", color, w.Count, true),
			newPlot("Median", color, w.Median, true),
			newPlot("95th Percentile", color, w.P95, true),
			newPlot("99th Percentile", color, w.P99, true),
		}
	case datagram.Set:
		s.plots = []*plot{
			newPlot("Unique Count", color, w.Unique, true),
			newPlot("Percent Unique", color, w.UniquePercent, true),
		}
	}

	s.row = termui.NewCol(dataWidth, 0, s.items()...)
	if prev != nil {
		s.row.Width = prev.Width
		s.row.Height = prev.Height
		s.row.X = prev.X
		s.row.Y = prev.Y
	}

	s.update()

	return s
}

func (s *plotSet) update() {
	ht := termui.TermHeight() - 1
	ct := len(s.plots)

	for _, p := range s.plots {
		p.update()

		p.lc.Height = ht / ct
		ht -= p.lc.Height
		ct--
	}

}

func (s *plotSet) items() []termui.GridBufferer {
	items := make([]termui.GridBufferer, 1+len(s.plots))
	items[0] = s.setHeader()

	for i, p := range s.plots {
		items[1+i] = p.chart()
	}

	return items
}

func (s *plotSet) setHeader() *termui.Par {
	p := termui.NewPar(s.metric.Name)

	p.Height = 1
	p.Border = false
	p.Bg = headerBorder
	p.TextBgColor = headerBorder
	p.TextFgColor = headerLabelColor

	return p
}
