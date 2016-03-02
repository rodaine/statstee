package views

import (
	"sync"

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
	db.Lock()
	db.next.draw()
	db.Unlock()
}

func (db *doubleBuffer) lazy() {
	go db.update(false)
}

func (db *doubleBuffer) update(force bool) {
	db.Lock()

	db.next.update(force)
	//db.current, db.next = db.next, db.current

	db.Unlock()
}
