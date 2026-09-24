package editor

import "unicode"

func (e *Editor) MoveUp() {
	if len(e.carets) > 0 {
		e.moveAll(e.MoveUp)
		return
	}
	if e.row > 0 {
		e.row--
		e.clampCol()
	}
	e.afterMove()
}

func (e *Editor) MoveDown() {
	if len(e.carets) > 0 {
		e.moveAll(e.MoveDown)
		return
	}
	if e.row < e.buf.LineCount()-1 {
		e.row++
		e.clampCol()
	}
	e.afterMove()
}

func (e *Editor) MoveLeft() {
	if len(e.carets) > 0 {
		e.moveAll(e.MoveLeft)
		return
	}
	if e.col > 0 {
		e.col--
	} else if e.row > 0 {
		e.row--
		e.col = e.buf.LineLen(e.row)
	}
	e.afterMove()
}

func (e *Editor) MoveRight() {
	if len(e.carets) > 0 {
		e.moveAll(e.MoveRight)
		return
	}
	if e.col < e.buf.LineLen(e.row) {
		e.col++
	} else if e.row < e.buf.LineCount()-1 {
		e.row++
		e.col = 0
	}
	e.afterMove()
}

func (e *Editor) LineStart() {
	if len(e.carets) > 0 {
		e.moveAll(e.LineStart)
		return
	}
	line := e.buf.Line(e.row)
	for i, r := range line {
		if !unicode.IsSpace(r) {
			e.col = i
			e.afterMove()
			return
		}
	}
	e.col = 0
	e.afterMove()
}

func (e *Editor) Home() {
	if len(e.carets) > 0 {
		e.moveAll(e.Home)
		return
	}
	e.col = 0
	e.afterMove()
}

func (e *Editor) MoveTo(row, col int) {
	e.row, e.col = row, col
	e.afterMove()
}

func (e *Editor) LineEnd() {
	if len(e.carets) > 0 {
		e.moveAll(e.LineEnd)
		return
	}
	e.col = e.buf.LineLen(e.row)
	e.afterMove()
}

func (e *Editor) PageUp() {
	e.row -= e.pageStep()
	if e.row < 0 {
		e.row = 0
	}
	e.clampCol()
	e.afterMove()
}

func (e *Editor) PageDown() {
	e.row += e.pageStep()
	if e.row > e.buf.LineCount()-1 {
		e.row = e.buf.LineCount() - 1
	}
	e.clampCol()
	e.afterMove()
}

func (e *Editor) GotoTop() {
	e.row, e.col = 0, 0
	e.afterMove()
}

func (e *Editor) GotoBottom() {
	e.row, e.col = e.buf.LineCount()-1, 0
	e.afterMove()
}

func (e *Editor) GotoLine(n int) {
	if n < 1 {
		n = 1
	}
	if n > e.buf.LineCount() {
		n = e.buf.LineCount()
	}
	e.row, e.col = n-1, 0
	e.afterMove()
}

func (e *Editor) WordNext() {
	for row := e.row; row < e.buf.LineCount(); row++ {
		line := e.buf.Line(row)
		start := 0
		if row == e.row {
			start = e.col + 1
		}
		for c := start; c < len(line); c++ {
			if unicode.IsSpace(line[c]) {
				continue
			}
			if c == 0 || isBoundary(line[c-1], line[c]) {
				e.row, e.col = row, c
				e.afterMove()
				return
			}
		}
	}
}

func (e *Editor) WordPrev() {
	for row := e.row; row >= 0; row-- {
		line := e.buf.Line(row)
		start := len(line) - 1
		if row == e.row {
			start = e.col - 1
		}
		for c := start; c > 0; c-- {
			if isBoundary(line[c-1], line[c]) {
				e.row, e.col = row, c
				e.afterMove()
				return
			}
		}
	}
}

func (e *Editor) pageStep() int {
	if e.rows > 2 {
		return e.rows / 2
	}
	return 1
}

func (e *Editor) clampCol() {
	if ll := e.buf.LineLen(e.row); e.col > ll {
		e.col = ll
	}
}

func (e *Editor) afterMove() {
	e.clampCursor()
	if e.sel.active {
		e.updateSelection()
	}
	e.scroll()
}

func isBoundary(prev, next rune) bool {
	return unicode.IsSpace(prev) && !unicode.IsSpace(next)
}
