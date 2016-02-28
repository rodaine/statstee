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
	router     *router.Router
	routerView *routerView
	dataView   *plotSet
}

func Display(r *router.Router) (err error) {
	if err = termui.Init(); err != nil {
		return
	}
	defer termui.Close()

	d := &display{
		router:     r,
		routerView: newRouterView(r),
		dataView:   newSet(bucket.DummyWindow),
	}

	termui.Body.AddRows(termui.NewRow(
		termui.NewCol(routerWidth, 0, d.routerView.list()),
		d.dataView.row,
	))

	d.update()

	d.registerHandlers()
	termui.Loop()
	return
}

func Quit() {
	termui.StopLoop()
}

func (d *display) registerHandlers() {
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		Quit()
	})

	termui.Handle("/sys/wnd/resize", func(e termui.Event) {
		d.forceUpdate()
		d.update()
	})

	termui.Handle("/sys/kbd/j", func(e termui.Event) {
		d.router.Next()
		d.update()
	})

	termui.Handle("/sys/kbd/k", func(e termui.Event) {
		d.router.Previous()
		d.update()
	})

	termui.Handle("/timer/1s", func(e termui.Event) {
		d.update()
	})
}

func (d *display) forceUpdate() {
	d.Lock()
	defer d.Unlock()

	termui.Body.Width = termui.TermWidth()
	termui.Body.Align()
}

func (d *display) update() {
	d.Lock()
	defer d.Unlock()

	defer func() {
		termui.Body.Align()
		termui.Render(termui.Body)
	}()

	d.routerView.update(d.router)

	if sel := d.router.SelectedMetric(); sel.Metric.Name != d.dataView.metric.Name {
		d.dataView = newSet(sel)
		termui.Body.Rows[0].Cols[1] = d.dataView.row
		termui.Body.Align()
	}
	d.dataView.update()
}
