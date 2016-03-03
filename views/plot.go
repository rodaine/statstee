package views

import (
	"fmt"
	"time"

	"github.com/gizak/termui"
	"github.com/rodaine/statstee/bucket"
)

const (
	plotLabelColor   = termui.ColorBlack | termui.AttrBold
	plotLineModifier = termui.AttrBold

	axesFormat = "15:04:05"
	valFormat  = "%s: %.4g"
	avgFormat  = "%s - avg(1m,5m,10m) %.4g, %.4g, %.4g"
)

type plotFunc func() []float64
type avgFunc func() bucket.Averages

type plot struct {
	avgs  avgFunc
	f     plotFunc
	label string
	lc    *termui.LineChart
}

func newPlot(label string, color termui.Attribute, f plotFunc, avgs avgFunc) *plot {
	lc := termui.NewLineChart()
	lc.BorderLabel = label
	lc.BorderBg = color

	p := &plot{
		f:     f,
		label: label,
		avgs:  avgs,
		lc:    lc,
	}
	p.reset()

	return p
}

func (p *plot) update() {
	p.reset()

	data := p.f()
	p.lc.Data = p.resize(data)
	p.lc.DataLabels = p.labels(len(p.lc.Data))

	p.lc.BorderLabel = fmt.Sprintf(
		valFormat,
		p.label,
		data[len(data)-1],
	)

	if p.avgs == nil {
		return
	}

	avgs := p.avgs()
	p.lc.BorderLabel = fmt.Sprintf(
		avgFormat,
		p.lc.BorderLabel,
		avgs.EWMA1,
		avgs.EWMA5,
		avgs.EWMA10,
	)
}

func (p *plot) chart() *termui.LineChart {
	return p.lc
}

func (p *plot) reset() {
	label := p.lc.BorderLabel
	color := p.lc.BorderBg

	lc := termui.NewLineChart()
	lc.Width = p.lc.Width
	lc.Height = p.lc.Height

	lc.BorderLeft = false
	lc.BorderRight = false
	lc.BorderBottom = false

	lc.BorderBg = color
	lc.BorderFg = color

	lc.BorderLabel = label
	lc.BorderLabelBg = color
	lc.BorderLabelFg = plotLabelColor

	lc.LineColor = color | plotLineModifier

	*p.lc = *lc
}

func (p *plot) resize(data []float64) []float64 {
	points := (p.lc.Width - 9) * 2
	offset := len(data) - points

	if offset <= 0 {
		return data
	}

	if offset < len(data) {
		return data[offset:]
	}

	return []float64{}
}

func (p *plot) labels(size int) []string {
	now := time.Now()

	lbls := make([]string, size)
	for i := 0; i < size; i++ {
		lbls[i] = now.Add(-1 * time.Duration(size-i) * time.Second).Format(axesFormat)
	}

	return lbls
}
