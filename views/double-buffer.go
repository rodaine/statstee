package views

import (
	"fmt"
	"log"
	"sync"

	"github.com/gizak/termui"
	"github.com/rodaine/statstee/router"
)

const (
	minWidth  = 78
	minHeight = 20
)

type doubleBuffer struct {
	sync.RWMutex
	width, height int
	current, next *display
}

func newBuffer(r *router.Router) *doubleBuffer {
	db := &doubleBuffer{
		current: newDisplay(r),
		next:    newDisplay(r),
	}

	db.dimSync()
	db.current.update()
	db.next.update()

	return db
}

func (db *doubleBuffer) draw() {
	var b termui.Bufferer
	db.RLock()
	if db.tooSmall() {
		b = db.smallView()
	} else {
		b = db.current.grid
	}

	termui.Render(b)
	db.RUnlock()
}

func (db *doubleBuffer) lazy() {
	go db.update(false)
}

func (db *doubleBuffer) update(force bool) {
	db.Lock()
	defer db.Unlock()

	if force {
		db.dimSync()
	}

	if db.tooSmall() {
		return
	}

	db.next.update()
	db.next.grid.Align()
	db.current, db.next = db.next, db.current
}

func (db *doubleBuffer) dimSync() {
	w, h := termui.TermWidth(), termui.TermHeight()

	db.width, db.height = w, h
	db.current.dimSync(w, h)
	db.next.dimSync(w, h)
}

func (db *doubleBuffer) smallView() *termui.Par {
	w := db.width
	h := db.height

	if w < 1 {
		w = 1
	}

	if h < 1 {
		h = 1
	}

	txt := fmt.Sprintf("SCREEN TOO SMALL\nCURRENT: %dx%d\nMIN: %dx%d", db.width, db.height, minWidth, minHeight)
	log.Println(txt)

	v := termui.NewPar(txt)
	v.Border = false
	v.Width = w
	v.Height = h

	v.Bg = termui.ColorRed
	v.TextBgColor = termui.ColorRed
	v.TextFgColor = termui.ColorBlack

	return v
}

func (db *doubleBuffer) tooSmall() bool {
	return db.width < minWidth || db.height < minHeight
}
