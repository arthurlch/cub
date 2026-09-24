package editor

const tabWidth = 4

func (e *Editor) InsertRune(r rune) {
	e.insert([]rune{r})
}

func (e *Editor) InsertNewline() {
	e.insert([]rune{'\n'})
}

func (e *Editor) InsertTab() {
	spaces := make([]rune, tabWidth)
	for i := range spaces {
		spaces[i] = ' '
	}
	e.insert(spaces)
}

func (e *Editor) insert(text []rune) {
	if len(e.carets) > 0 {
		e.editAll(func() { e.recordInsert(text) })
		return
	}
	e.recordInsert(text)
	e.scroll()
}

func (e *Editor) Backspace() {
	if e.sel.active {
		e.DeleteSelection()
		return
	}
	do := func() {
		if e.col > 0 {
			e.recordDelete(e.row, e.col-1, e.row, e.col)
		} else if e.row > 0 {
			prev := e.buf.LineLen(e.row - 1)
			e.recordDelete(e.row-1, prev, e.row, 0)
		}
	}
	if len(e.carets) > 0 {
		e.editAll(do)
		return
	}
	do()
	e.scroll()
}

func (e *Editor) DeleteForward() {
	if e.sel.active {
		e.DeleteSelection()
		return
	}
	do := func() {
		ll := e.buf.LineLen(e.row)
		if e.col < ll {
			e.recordDelete(e.row, e.col, e.row, e.col+1)
		} else if e.row < e.buf.LineCount()-1 {
			e.recordDelete(e.row, e.col, e.row+1, 0)
		}
	}
	if len(e.carets) > 0 {
		e.editAll(do)
		return
	}
	do()
	e.scroll()
}

func (e *Editor) Append() {
	if e.col < e.buf.LineLen(e.row) {
		e.col++
	}
	e.EnterInsert()
}

func (e *Editor) OpenBelow() {
	e.col = e.buf.LineLen(e.row)
	e.recordInsert([]rune{'\n'})
	e.mode = InsertMode
	e.scroll()
}

func (e *Editor) OpenAbove() {
	e.col = 0
	e.recordInsert([]rune{'\n'})
	if e.row > 0 {
		e.row--
	}
	e.col = 0
	e.mode = InsertMode
	e.scroll()
}

func (e *Editor) DeleteChar() {
	if e.sel.active {
		e.DeleteSelection()
		return
	}
	if e.col < e.buf.LineLen(e.row) {
		e.recordDelete(e.row, e.col, e.row, e.col+1)
	}
	e.clampCursor()
	e.scroll()
}

func (e *Editor) DeleteLine() {
	e.yankLine()
	last := e.buf.LineCount() - 1
	if e.row < last {
		e.recordDelete(e.row, 0, e.row+1, 0)
	} else if e.row > 0 {
		prev := e.buf.LineLen(e.row - 1)
		e.recordDelete(e.row-1, prev, e.row, e.buf.LineLen(e.row))
	} else {
		e.recordDelete(e.row, 0, e.row, e.buf.LineLen(e.row))
	}
	e.clampCursor()
	e.scroll()
}

func (e *Editor) YankLine() { e.yankLine() }

func (e *Editor) yankLine() {
	line := append([]rune{}, e.buf.Line(e.row)...)
	e.clip = append(line, '\n')
	e.clipLinewise = true
}
