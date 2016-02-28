package views

import (
	"math"
	"time"

	"fmt"

	"github.com/gizak/termui"
)

const (
	plotLabelColor   = termui.ColorBlack | termui.AttrBold
	plotLineModifier = termui.AttrBold

	axesFormat = "15:04:05"
	valFormat  = "%s: %.3g"
	avgFormat  = "%s - avg(1m,5m,10m) %.2g, %.2g, %.2g"
)

type plotFunc func() []float64

type plot struct {
	avgs  bool
	f     plotFunc
	label string
	lc    *termui.LineChart
}

func newPlot(label string, color termui.Attribute, f plotFunc, avgs bool) *plot {
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

	if !p.avgs {
		return
	}

	p.lc.BorderLabel = fmt.Sprintf(
		avgFormat,
		p.lc.BorderLabel,
		ewma(data, 1),
		ewma(data, 5),
		ewma(data, 10),
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

// based off: https://en.wikipedia.org/wiki/Moving_average#Application_to_measuring_computer_performance
func ewma(data []float64, minutes float64) float64 {
	ct := len(data)
	W := float64(ct) / 60.0 // How long is the data relevant (in minutes)
	a := math.Exp(-1.0 / (W * minutes))

	y := data[0]
	for i := 1; i < ct; i++ {
		y = data[i] + a*(y-data[i])
	}

	return y
}
