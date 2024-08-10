package state

/*

This file is holding the state for the whole text editor
State that is responsible for the buffer is the most important
Using rune has data structure so we can represent multi dimension word, lines, files.
Grid is composed of Rows and Columns.
Current mode for the editor are view and insert modes, choice was made to use iota.

*/

import (
	"time"

	"github.com/nsf/termbox-go"
)

type Mode int

const (
	ViewMode Mode = iota
	InsertMode
)

type ChangeType int

const (
	Insert ChangeType = iota
	Delete
)

type Change struct {
	Type ChangeType
	Row, Col  int
	Text      []rune
	PrevRow, PrevCol int
}


type State struct {
	Mode        Mode
	Rows, Cols  int
	OffsetRow   int
	OffsetCol   int
	CurrentRow  int
	CurrentCol  int
	StartRow    int
	StartCol    int
	EndRow      int
	EndCol 			int
	SourceFile  string
	SelectionActive bool
	TextBuffer  [][]rune
	UndoBuffer  []Change
  RedoBuffer  []Change
	CopyBuffer  []rune
	Modified    bool
	QuitKey     termbox.Key
	ErrorMessage string
	MessageTimestamp time.Time
	LastKey rune
}
