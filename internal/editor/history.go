package editor

// THE HACKIEST HISTORY SYSTEM EVER !!!
// I f hate it
const maxHistoryDepth = 1000

type command struct {
	insert         bool
	row, col       int
	endRow, endCol int
	text           []rune

	beforeRow, beforeCol int
	afterRow, afterCol   int
}

type history struct {
	undo [][]command
	redo [][]command
}

func (e *Editor) recordInsert(text []rune) {
	br, bc := e.row, e.col
	er, ec := e.buf.Insert(br, bc, text)
	e.row, e.col = er, ec
	e.push(command{
		insert: true, row: br, col: bc, endRow: er, endCol: ec,
		text:      append([]rune{}, text...),
		beforeRow: br, beforeCol: bc, afterRow: er, afterCol: ec,
	})
	e.modified = true
}

func (e *Editor) recordDelete(sr, sc, er, ec int) {
	sr, sc, er, ec = order(sr, sc, er, ec)
	if sr == er && sc == ec {
		return
	}
	br, bc := e.row, e.col
	removed := e.buf.Delete(sr, sc, er, ec)
	e.row, e.col = sr, sc
	e.push(command{
		insert: false, row: sr, col: sc, endRow: er, endCol: ec,
		text:      removed,
		beforeRow: br, beforeCol: bc, afterRow: sr, afterCol: sc,
	})
	e.modified = true
}

func (e *Editor) push(c command) {
	if e.grouping {
		e.group = append(e.group, c)
		return
	}
	e.pushGroup([]command{c})
}

func (e *Editor) pushGroup(g []command) {
	if len(g) == 0 {
		return
	}
	e.hist.undo = append(e.hist.undo, g)
	if len(e.hist.undo) > maxHistoryDepth {
		e.hist.undo = e.hist.undo[len(e.hist.undo)-maxHistoryDepth:]
	}
	e.hist.redo = nil
}

func (e *Editor) beginGroup() {
	e.grouping = true
	e.group = nil
}

func (e *Editor) commitGroup() {
	e.grouping = false
	g := e.group
	e.group = nil
	e.pushGroup(g)
}

func (e *Editor) Undo() {
	n := len(e.hist.undo)
	if n == 0 {
		return
	}
	g := e.hist.undo[n-1]
	e.hist.undo = e.hist.undo[:n-1]
	for i := len(g) - 1; i >= 0; i-- {
		c := g[i]
		if c.insert {
			e.buf.Delete(c.row, c.col, c.endRow, c.endCol)
		} else {
			e.buf.Insert(c.row, c.col, c.text)
		}
	}
	e.carets = nil
	e.row, e.col = g[0].beforeRow, g[0].beforeCol
	e.hist.redo = append(e.hist.redo, g)
	e.modified = true
	e.clampCursor()
}

func (e *Editor) Redo() {
	n := len(e.hist.redo)
	if n == 0 {
		return
	}
	g := e.hist.redo[n-1]
	e.hist.redo = e.hist.redo[:n-1]
	for _, c := range g {
		if c.insert {
			e.buf.Insert(c.row, c.col, c.text)
		} else {
			e.buf.Delete(c.row, c.col, c.endRow, c.endCol)
		}
	}
	e.carets = nil
	e.row, e.col = g[len(g)-1].afterRow, g[len(g)-1].afterCol
	e.hist.undo = append(e.hist.undo, g)
	e.modified = true
	e.clampCursor()
}

func (e *Editor) CanUndo() bool { return len(e.hist.undo) > 0 }
func (e *Editor) CanRedo() bool { return len(e.hist.redo) > 0 }

func order(sr, sc, er, ec int) (int, int, int, int) {
	if sr > er || (sr == er && sc > ec) {
		return er, ec, sr, sc
	}
	return sr, sc, er, ec
}
