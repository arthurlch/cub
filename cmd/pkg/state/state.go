package state

import (
	"time"

	"github.com/nsf/termbox-go"
)

type Mode int

const (
	ViewMode Mode = iota
	InsertMode
)

type Buffer struct {
	TextBuffer [][]rune
	SourceFile string
	Modified   bool
	Language   string
}

type Cursor struct {
	CurrentRow int
	CurrentCol int
}

type Selection struct {
	StartRow        int
	StartCol        int
	EndRow          int
	EndCol          int
	SelectionActive bool
}

type Viewport struct {
	OffsetRow int
	OffsetCol int
	Rows      int
	Cols      int
}

type UIState struct {
	ErrorMessage     string
	MessageTimestamp time.Time
	LastKey          rune
	QuitKey          termbox.Key
	LineNumberBuffer string
}

type State struct {
	Mode Mode
	Buffer
	Cursor
	Selection
	Viewport

	UIState

	CopyBuffer []rune

	Quit bool

	history History
}
