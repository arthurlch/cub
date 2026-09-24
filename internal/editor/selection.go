package editor

// who said go became verbose ?? !

func (e *Editor) StartSelection() {
	e.sel = selection{active: true, startRow: e.row, startCol: e.col, endRow: e.row, endCol: e.col}
}

func (e *Editor) EndSelection() { e.endSelection() }

func (e *Editor) SelectionActive() bool { return e.sel.active }

func (e *Editor) endSelection() { e.sel.active = false }

func (e *Editor) updateSelection() {
	e.sel.endRow, e.sel.endCol = e.row, e.col
}

func (e *Editor) SelectAll() {
	if e.buf.LineCount() == 0 {
		return
	}
	last := e.buf.LineCount() - 1
	e.sel = selection{
		active:   true,
		startRow: 0, startCol: 0,
		endRow: last, endCol: e.buf.LineLen(last),
	}
	e.row, e.col = last, e.buf.LineLen(last)
	e.afterMove()
}

func (e *Editor) selectedText() []rune {
	if !e.sel.active {
		return nil
	}
	sr, sc, er, ec := order(e.sel.startRow, e.sel.startCol, e.sel.endRow, e.sel.endCol)
	sc = clampInt(sc, 0, e.buf.LineLen(sr))
	ec = clampInt(ec, 0, e.buf.LineLen(er))

	if sr == er {
		return append([]rune{}, e.buf.Line(sr)[sc:ec]...)
	}
	out := append([]rune{}, e.buf.Line(sr)[sc:]...)
	for r := sr + 1; r < er; r++ {
		out = append(out, '\n')
		out = append(out, e.buf.Line(r)...)
	}
	out = append(out, '\n')
	out = append(out, e.buf.Line(er)[:ec]...)
	return out
}

func (e *Editor) Copy() {
	if text := e.selectedText(); text != nil {
		e.clip = text
		e.clipLinewise = false
	}
}

func (e *Editor) Cut() {
	if !e.sel.active {
		return
	}
	e.clip = e.selectedText()
	e.clipLinewise = false
	e.DeleteSelection()
}

func (e *Editor) DeleteSelection() {
	if !e.sel.active {
		return
	}
	sr, sc, er, ec := order(e.sel.startRow, e.sel.startCol, e.sel.endRow, e.sel.endCol)
	e.endSelection()
	e.recordDelete(sr, sc, er, ec)
	e.scroll()
}

func (e *Editor) Paste() {
	if len(e.clip) == 0 {
		return
	}
	if e.sel.active {
		e.DeleteSelection()
	}
	if e.clipLinewise {
		e.col = e.buf.LineLen(e.row)
		e.recordInsert(append([]rune{'\n'}, e.lineClip()...))
	} else {
		e.recordInsert(e.clip)
	}
	e.scroll()
}

func (e *Editor) PasteBefore() {
	if len(e.clip) == 0 {
		return
	}
	if e.sel.active {
		e.DeleteSelection()
	}
	if e.clipLinewise {
		e.col = 0
		e.recordInsert(append(e.lineClip(), '\n'))
		if e.row > 0 {
			e.row--
		}
		e.col = 0
	} else {
		e.recordInsert(e.clip)
	}
	e.scroll()
}

func (e *Editor) lineClip() []rune {
	if n := len(e.clip); n > 0 && e.clip[n-1] == '\n' {
		return append([]rune{}, e.clip[:n-1]...)
	}
	return append([]rune{}, e.clip...)
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
