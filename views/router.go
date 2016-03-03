package views

import (
	"fmt"
	"strings"

	"github.com/gizak/termui"
	"github.com/rodaine/statstee/datagram"
	"github.com/rodaine/statstee/router"
)

const (
	instructions     = "q:Quit j:Next k:Prev"
	headerTextFormat = "%s | %d/%d"
	headerTextColor  = termui.ColorBlack | termui.AttrBold

	borderColor = termui.ColorWhite

	unselectedFormat = " [%s] %s"
	selectedFormat   = "[%s](%s,BOLD)"

	selectedPadding = " "
	selectedColor   = "black"
)

type routerView struct {
	offset int
	l      *termui.List
	height int
}

func newRouterView(r *router.Router) *routerView {
	v := &routerView{l: termui.NewList()}

	v.l.BorderLabel = instructions
	v.l.BorderLabelFg = headerTextColor
	v.l.BorderLabelBg = borderColor

	v.l.BorderFg = borderColor
	v.l.BorderBg = borderColor

	return v
}

func (v *routerView) update(r *router.Router) {
	s := r.Selected()
	ms := r.Metrics()

	v.l.Height = v.height

	items := make([]string, len(ms))
	selIdx := 0

	for i, m := range ms {
		selected := m.Name == s
		items[i] = v.metricItemLabel(m, selected)
		if selected {
			selIdx = i
		}
	}

	v.l.BorderLabel = fmt.Sprintf(headerTextFormat, instructions, selIdx+1, len(ms))

	max := v.offset + v.l.Height - 3
	if selIdx >= v.offset && selIdx <= max {
		// noop
	} else if selIdx < v.offset {
		v.offset = selIdx
	} else {
		v.offset += selIdx - max
	}

	v.l.Items = items[v.offset:]
}

func (v *routerView) metricItemLabel(m datagram.Metric, selected bool) string {
	lbl := fmt.Sprintf(unselectedFormat, m.TypePrefix(), m.Name)
	if !selected {
		return lbl
	}

	offset := v.l.Width - len(lbl)
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
	return v.l
}
