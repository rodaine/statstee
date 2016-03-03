package views

import (
	"sync"

	"github.com/gizak/termui"
	"github.com/rodaine/statstee/bucket"
	"github.com/rodaine/statstee/router"
)

const (
	maxGridSpan    = 12
	routerGridSpan = 3
	dataGridSpan   = maxGridSpan - routerGridSpan
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
		routerView: newRouterView(r),
		dataView:   newSet(bucket.DummyWindow, nil),
	}

	d.grid = termui.NewGrid(termui.NewRow(
		termui.NewCol(routerGridSpan, 0, d.routerView.list()),
		d.dataView.row,
	))

	return d
}

func (d *display) update() {
	d.routerView.update(d.router)

	if sel := d.router.SelectedMetric(); sel.Metric.Name != d.dataView.metric.Name {
		d.dataView = newSet(sel, d.grid.Rows[0].Cols[1])
		d.dataView.height = d.height
		d.grid.Rows[0].Cols[1] = d.dataView.row
		d.grid.Align()
	}

	d.dataView.update()
}

func (d *display) dimSync(width, height int) {
	d.width, d.height = width, height

	d.routerView.height = height
	d.dataView.height = height

	d.grid.Width = width
	d.grid.Align()
}
