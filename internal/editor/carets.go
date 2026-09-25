package editor

import "sort"

type caret struct {
	row, col int
}

func (e *Editor) CaretCount() int { return len(e.carets) + 1 }

func (e *Editor) Carets() [][2]int {
	out := make([][2]int, 0, len(e.carets)+1)
	out = append(out, [2]int{e.row, e.col})
	for _, c := range e.carets {
		out = append(out, [2]int{c.row, c.col})
	}
	return out
}

func (e *Editor) ClearCarets() { e.carets = nil }

func (e *Editor) AddCaretBelow() {
	row := e.maxCaretRow() + 1
	if row >= e.buf.LineCount() {
		return
	}
	e.addCaret(row, e.col)
}

func (e *Editor) AddCaretAbove() {
	row := e.minCaretRow() - 1
	if row < 0 {
		return
	}
	e.addCaret(row, e.col)
}

func (e *Editor) addCaret(row, col int) {
	if ll := e.buf.LineLen(row); col > ll {
		col = ll
	}
	if row == e.row && col == e.col {
		return
	}
	for _, c := range e.carets {
		if c.row == row && c.col == col {
			return
		}
	}
	e.carets = append(e.carets, caret{row, col})
}

func (e *Editor) maxCaretRow() int {
	m := e.row
	for _, c := range e.carets {
		if c.row > m {
			m = c.row
		}
	}
	return m
}

func (e *Editor) minCaretRow() int {
	m := e.row
	for _, c := range e.carets {
		if c.row < m {
			m = c.row
		}
	}
	return m
}

func (e *Editor) allCarets() []caret {
	out := make([]caret, 0, len(e.carets)+1)
	out = append(out, caret{e.row, e.col})
	out = append(out, e.carets...)
	return out
}

func (e *Editor) setResults(out []caret) {
	out = dedupeCarets(out)
	e.row, e.col = out[0].row, out[0].col
	e.carets = out[1:]
}

func (e *Editor) moveAll(do func()) {
	pts := e.allCarets()
	e.carets = nil
	out := make([]caret, 0, len(pts))
	for _, p := range pts {
		e.row, e.col = p.row, p.col
		do()
		out = append(out, caret{e.row, e.col})
	}
	e.setResults(out)
	e.scroll()
}

func (e *Editor) editAll(do func()) {
	pts := e.allCarets()
	e.carets = nil
	order := descByPosition(pts)
	results := make([]caret, len(pts))
	done := make([]int, 0, len(pts))

	e.beginGroup()
	for _, i := range order {
		before := e.buf.LineCount()
		e.row, e.col = pts[i].row, pts[i].col
		do()
		results[i] = caret{e.row, e.col}
		if delta := e.buf.LineCount() - before; delta != 0 {
			for _, j := range done {
				if results[j].row > pts[i].row {
					results[j].row += delta
				}
			}
		}
		done = append(done, i)
	}
	e.commitGroup()

	e.setResults(results)
	e.scroll()
}

func descByPosition(pts []caret) []int {
	idx := make([]int, len(pts))
	for i := range idx {
		idx[i] = i
	}
	sort.SliceStable(idx, func(a, b int) bool {
		x, y := pts[idx[a]], pts[idx[b]]
		if x.row != y.row {
			return x.row > y.row
		}
		return x.col > y.col
	})
	return idx
}

func dedupeCarets(in []caret) []caret {
	seen := make(map[caret]bool, len(in))
	var out []caret
	for _, c := range in {
		if seen[c] {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	return out
}
