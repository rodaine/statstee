package views

import (
	"fmt"
	"strings"

	"log"

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
	offset      int
	l           *termui.List
	needsRedraw bool
}

func newRouterView(r *router.Router) *routerView {
	v := &routerView{l: termui.NewList()}

	v.l.BorderLabel = instructions
	v.l.BorderLabelFg = headerTextColor
	v.l.BorderLabelBg = borderColor

	v.l.BorderFg = borderColor
	v.l.BorderBg = borderColor

	v.update(r, true)
	return v
}

func (v *routerView) drawIfNeeded() {
	if !v.needsRedraw {
		return
	}
	v.needsRedraw = false

	log.Println("drawing router")
	termui.Render(v.l)
}

func (v *routerView) update(r *router.Router, force bool) {
	if !r.NeedsUpdate() && !force {
		return
	}
	v.needsRedraw = true

	s := r.Selected()
	ms := r.Metrics()

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
	v.l.Height = termui.TermHeight()

	max := v.offset + v.l.Height - 3
	if selIdx >= v.offset && selIdx <= max {
		// noop
	} else if selIdx < v.offset {
		v.offset = selIdx
	} else {
		v.offset += selIdx - max
	}

	if len(items) > 0 {
		items = items[v.offset:]
	}

	v.l.Items = items
	log.Println("updated router view")
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
