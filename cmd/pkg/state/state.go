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

type UndoState struct {
	TextBuffer [][]rune
	CurrentRow int
	CurrentCol int
}
	
type State struct {
	Mode             Mode
	Rows, Cols       int
	OffsetRow        int
	OffsetCol        int
	CurrentRow       int
	CurrentCol       int
	StartRow         int
	StartCol         int
	EndRow           int
	EndCol           int
	SourceFile       string
	SelectionActive  bool
	TextBuffer       [][]rune
	UndoBuffer       []UndoState
	RedoBuffer       []UndoState
	CopyBuffer       []rune
	Modified         bool
	QuitKey          termbox.Key
	ErrorMessage     string
	MessageTimestamp time.Time
	LastKey          rune
	LineNumberBuffer string
}
