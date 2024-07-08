package state

import (
	"time"

	"github.com/nsf/termbox-go"
)

type Mode int

const (
	ViewMode Mode = iota
	EditMode
)


type State struct {
	Mode        Mode
	Rows, Cols  int
	OffsetRow   int
	OffsetCol   int
	CurrentRow  int
	CurrentCol  int
	StartRow       int
	StartCol       int
	SourceFile  string
	TextBuffer  [][]rune
	UndoBuffer  [][]rune
	CopyBuffer  []rune
	Modified    bool
	QuitKey     termbox.Key
	ErrorMessage string
	MessageTimestamp time.Time
}
