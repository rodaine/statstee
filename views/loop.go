package views

import (
	"github.com/gizak/termui"
	"github.com/rodaine/statstee/router"
)

func Loop(r *router.Router) (err error) {
	if err = termui.Init(); err != nil {
		return
	}
	defer termui.Close()
	defer recover()

	db := newBuffer(r)
	db.draw()

	registerHandlers(r, db)
	termui.Loop()
	return
}

func Quit() {
	termui.StopLoop()
}

func registerHandlers(r *router.Router, db *doubleBuffer) {
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		Quit()
	})

	termui.Handle("/sys/wnd/resize", func(e termui.Event) {
		db.update(true)
		db.draw()
	})

	termui.Handle("/sys/kbd/j", func(e termui.Event) {
		r.Next()
		db.update(false)
		db.draw()
	})

	termui.Handle("/sys/kbd/k", func(e termui.Event) {
		r.Previous()
		db.update(false)
		db.draw()
	})

	termui.Handle("/timer/1s", func(e termui.Event) {
		db.draw()
		db.lazy()
	})
}
