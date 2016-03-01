package views

import (
	"fmt"
	"strings"

	"github.com/gizak/termui"
	"github.com/rodaine/statstee/datagram"
	"github.com/rodaine/statstee/router"
)

const (
	headerText      = "q:Quit j:Next k:Prev"
	headerTextColor = termui.ColorBlack | termui.AttrBold

	borderColor = termui.ColorWhite

	unselectedFormat = " [%s] %s"
	selectedFormat   = "[%s](%s,BOLD)"

	selectedPadding = " "
	selectedColor   = "black"
)

type routerView termui.List

func newRouterView(r *router.Router) *routerView {
	v := (*routerView)(termui.NewList())

	v.BorderLabel = headerText
	v.BorderLabelFg = headerTextColor
	v.BorderLabelBg = borderColor

	v.BorderFg = borderColor
	v.BorderBg = borderColor

	v.update(r)
	return v
}

func (v *routerView) update(r *router.Router) {
	s := r.Selected()
	ms := r.Metrics()

	v.Items = make([]string, len(ms))
	for i, m := range ms {
		v.Items[i] = v.metricItemLabel(m, m.Name == s)
	}

	v.Height = termui.TermHeight()
}

func (v *routerView) metricItemLabel(m datagram.Metric, selected bool) string {
	lbl := fmt.Sprintf(unselectedFormat, m.TypePrefix(), m.Name)
	if !selected {
		return lbl
	}

	offset := v.Width - len(lbl)
	if offset > 0 {
		lbl += strings.Repeat(selectedPadding, offset)
	}

	return fmt.Sprintf(
		selectedFormat,
		lbl,
		markdown(selectedColor, datagramColorNames[m.Type]),
	)
}

func (v *routerView) list() *termui.List {
	return (*termui.List)(v)
}
