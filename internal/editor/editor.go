package editor

import (
	"time"

	"github.com/arthurlch/cub/internal/buffer"
	"github.com/arthurlch/cub/internal/document"
)

type Mode int

const (
	ViewMode Mode = iota
	InsertMode
)

const messageTTL = 4 * time.Second

type selection struct {
	active             bool
	startRow, startCol int
	endRow, endCol     int
}

type Editor struct {
	buf  *buffer.Buffer
	mode Mode

	row, col       int
	offRow, offCol int
	rows, cols     int

	sel          selection
	clip         []rune
	clipLinewise bool
	hist         history

	carets   []caret
	grouping bool
	group    []command

	file     string
	language string
	modified bool

	message   string
	messageAt time.Time
}

func New() *Editor {
	return &Editor{buf: buffer.New(""), language: "Plain Text"}
}

func (e *Editor) SetText(s string) {
	e.buf.SetText(s)
	e.row, e.col = 0, 0
	e.offRow, e.offCol = 0, 0
	e.endSelection()
	e.carets = nil
	e.hist = history{}
	e.modified = false
}

func (e *Editor) Buffer() *buffer.Buffer { return e.buf }
func (e *Editor) LineCount() int         { return e.buf.LineCount() }
func (e *Editor) Line(i int) []rune      { return e.buf.Line(i) }
func (e *Editor) Cursor() (int, int)     { return e.row, e.col }
func (e *Editor) Offset() (int, int)     { return e.offRow, e.offCol }
func (e *Editor) Mode() Mode             { return e.mode }
func (e *Editor) FileName() string       { return e.file }
func (e *Editor) Language() string       { return e.language }
func (e *Editor) Modified() bool         { return e.modified }

func (e *Editor) Selection() (active bool, sr, sc, er, ec int) {
	s := e.sel
	return s.active, s.startRow, s.startCol, s.endRow, s.endCol
}

func (e *Editor) EnterInsert() {
	e.mode = InsertMode
	e.endSelection()
	e.clampCursor()
}

func (e *Editor) EnterView() { e.mode = ViewMode }

func (e *Editor) SetViewport(rows, cols int) {
	e.rows, e.cols = rows, cols
	e.scroll()
}

func (e *Editor) ShowMessage(msg string) {
	e.message = msg
	e.messageAt = time.Now()
}

func (e *Editor) ActiveMessage() (string, bool) {
	if e.message == "" || time.Since(e.messageAt) > messageTTL {
		return "", false
	}
	return e.message, true
}

func (e *Editor) clampCursor() {
	if e.row < 0 {
		e.row = 0
	}
	if e.row >= e.buf.LineCount() {
		e.row = e.buf.LineCount() - 1
	}
	ll := e.buf.LineLen(e.row)
	if e.col < 0 {
		e.col = 0
	}
	if e.col > ll {
		e.col = ll
	}
}

func (e *Editor) scroll() {
	if e.rows <= 0 || e.cols <= 0 {
		return
	}
	if e.row < e.offRow {
		e.offRow = e.row
	}
	if e.row >= e.offRow+e.rows {
		e.offRow = e.row - e.rows + 1
	}
	if e.col < e.offCol {
		e.offCol = e.col
	}
	if e.col >= e.offCol+e.cols {
		e.offCol = e.col - e.cols + 1
	}
	if e.offRow < 0 {
		e.offRow = 0
	}
	if e.offCol < 0 {
		e.offCol = 0
	}
}

func (e *Editor) refreshLanguage() {
	e.language = document.Language(e.file)
}
