package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type EditorState struct {
	State *state.State
}

func NewEditorState(sharedState *state.State) *EditorState {
	return &EditorState{State: sharedState}
}

func (es *EditorState) StatusBar() {
	st := es.State
	filename := filepath.Base(st.SourceFile)
	if len(filename) > 14 {
		filename = filename[:14]
	}
	
	leftStatus := filename
	if st.Modified {
		leftStatus += " [MODIFIED]"
	} else {
		leftStatus += " [SAVED]"
	}

	modeStatus := " VIEW "
	if st.Mode == state.InsertMode {
		modeStatus = " INSERT "
	}
	leftStatus = modeStatus + leftStatus

	rightStatus := fmt.Sprintf("Row %d Col %d ", st.CurrentRow+1, st.CurrentCol)
	
	padding := st.Cols - len(leftStatus) - len(rightStatus)
	if padding < 0 {
		padding = 0
	}

	fullStatusBar := leftStatus + strings.Repeat(" ", padding) + rightStatus

	printMessage(0, st.Rows, StatusBarForeground, StatusBarBackground, fullStatusBar)
}

func printMessage(col, row int, foreground, background termbox.Attribute, message string) {
	for _, ch := range message {
		termbox.SetCell(col, row, ch, foreground, background)
		col += runewidth.RuneWidth(ch)
	}
}