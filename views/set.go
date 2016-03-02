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
	height int
}

func newSet(w *bucket.MetricWindow, prev *termui.Row, height int) *plotSet {
	s := &plotSet{metric: w.Metric, height: height}
	s.header = s.setHeader()

	color := datagramColors[w.Metric.Type]

	switch w.Metric.Type {
	case datagram.Gauge:
		s.plots = []*plot{
			newPlot("Gauge Value", color, w.Last, w.LastAverages),
		}
	case datagram.Counter:
		s.plots = []*plot{
			newPlot("Count", color, w.Sum, w.SumAverages),
			newPlot("Cumulative Count", color, w.CumSum, nil),
		}
	case datagram.Histogram, datagram.Timer:
		s.plots = []*plot{
			newPlot("Count", color, w.Count, w.CountAverages),
			newPlot("Median", color, w.Median, w.MedianAverages),
			newPlot("75th Percentile", color, w.P75, w.P75Averages),
			newPlot("95th Percentile", color, w.P75, w.P95Averages),
		}
	case datagram.Set:
		s.plots = []*plot{
			newPlot("Unique Count", color, w.Unique, w.UniqueAverages),
			newPlot("Percent Unique", color, w.UniquePercent, w.UniquePercentAverages),
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

func (s *plotSet) draw() {
	termui.Render(s.row)
}

func (s *plotSet) update() {
	ht := s.height - 1
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
