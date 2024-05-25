package ui

import (
	"fmt"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type EditorState struct {
	State *state.State
}

func NewEditorState() *EditorState {
	return &EditorState{State: &state.State{}}
}

func (es *EditorState) StatusBar() {
	st := es.State
	filename := st.SourceFile
	if len(filename) > 14 {
		filename = filename[:14]
	}
	fileStatus := fmt.Sprintf("%s - %d lines", filename, len(st.TextBuffer))
	if st.Modified {
		fileStatus += " modified"
	} else {
		fileStatus += " saved"
	}

	modeStatus := "VIEW: "
	if st.Mode == state.EditMode {
		modeStatus = "EDIT: "
	}

	cursorStatus := fmt.Sprintf("Row %d Col %d ", st.CurrentRow+1, st.CurrentCol)
	statusBar := modeStatus + fileStatus + cursorStatus
	termbox.SetCursor(0, st.Rows)
	print_message(0, st.Rows, termbox.ColorBlack, termbox.ColorWhite, statusBar)
}

func print_message(col, row int, foreground, background termbox.Attribute, message string) {
	for _, ch := range message {
		termbox.SetCell(col, row, ch, foreground, background)
		col += runewidth.RuneWidth(ch)
	}
}
