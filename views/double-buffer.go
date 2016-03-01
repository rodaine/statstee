package views

import (
	"sync"

	"log"

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
	log.Println("lazy")
	go db.update(false)
}

func (db *doubleBuffer) update(force bool) {
	db.Lock()
	log.Println("updating")
	db.next.update(false)
	db.next.grid.Width = termui.TermWidth()
	db.next.grid.Align()
	db.current, db.next = db.next, db.current

	log.Println("updated")
	db.Unlock()
}
