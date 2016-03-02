package views

import (
	"sync"

	"github.com/gizak/termui"
	"github.com/rodaine/statstee/bucket"
	"github.com/rodaine/statstee/router"
)

const (
	maxWidth    = 12
	routerWidth = 3
	dataWidth   = maxWidth - routerWidth
)

type display struct {
	sync.RWMutex
	router *router.Router

	grid       *termui.Grid
	routerView *routerView
	dataView   *plotSet

	width, height int
}

func newDisplay(r *router.Router) *display {
	d := &display{
		router:     r,
		routerView: newRouterView(r, 100),
		dataView:   newSet(bucket.DummyWindow, nil, 100),
	}

	d.grid = termui.NewGrid(termui.NewRow(
		termui.NewCol(routerWidth, 0, d.routerView.l),
		d.dataView.row,
	))

	d.refreshDims()

	return d
}

func (d *display) draw() {
	termui.Render(d.grid)
}

func (d *display) update(force bool) {
	if force {
		d.refreshDims()
	}

	d.routerView.update(d.router, force)

	if sel := d.router.SelectedMetric(); sel.Metric.Name != d.dataView.metric.Name {
		d.dataView = newSet(sel, d.grid.Rows[0].Cols[1], d.height)
		d.grid.Rows[0].Cols[1] = d.dataView.row
		d.grid.Align()
	}

	d.dataView.update()
}

func (d *display) refreshDims() {
	d.width, d.height = termui.TermWidth(), termui.TermHeight()

	d.grid.Width = d.width

	d.routerView.height = d.height
	d.dataView.height = d.height

	d.grid.Align()
}
