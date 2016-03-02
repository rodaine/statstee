package views

import (
	"sync"

	"github.com/gizak/termui"
	"github.com/rodaine/statstee/router"
)

type doubleBuffer struct {
	sync.RWMutex
	current, next *display
}

func newBuffer(r *router.Router) *doubleBuffer {
	return &doubleBuffer{
		current: newDisplay(r),
		next:    newDisplay(r),
	}
}

func (db *doubleBuffer) draw() {
	db.RLock()
	g := db.current.grid
	db.RUnlock()

	termui.Render(g)
}

func (db *doubleBuffer) lazy() {
	go db.update(false)
}

func (db *doubleBuffer) update(force bool) {
	db.Lock()

	db.next.update(false)
	db.next.grid.Width = termui.TermWidth()
	db.next.grid.Align()
	db.current, db.next = db.next, db.current

	db.Unlock()
}
