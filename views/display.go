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
}

func newDisplay(r *router.Router) *display {
	d := &display{
		router:     r,
		routerView: newRouterView(r),
		dataView:   newSet(bucket.DummyWindow, nil),
	}

	d.grid = termui.NewGrid(termui.NewRow(
		termui.NewCol(routerWidth, 0, d.routerView.list()),
		d.dataView.row,
	))
	d.grid.Width = termui.TermWidth()
	d.grid.Align()

	return d
}

func (d *display) update(force bool) {
	if force {
		d.grid.Width = termui.TermWidth()
		d.grid.Align()
	}

	d.routerView.update(d.router)

	if sel := d.router.SelectedMetric(); sel.Metric.Name != d.dataView.metric.Name {
		d.dataView = newSet(sel, d.grid.Rows[0].Cols[1])
		d.grid.Rows[0].Cols[1] = d.dataView.row
		d.grid.Align()
	}

	d.dataView.update()
}
