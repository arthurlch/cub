package render

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/theme"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func StatusBar(st *state.State) {
	filename := filepath.Base(st.SourceFile)
	if len(filename) > 14 {
		filename = filename[:14]
	}

	fileStatus := filename
	if st.Modified {
		fileStatus += " [MODIFIED]"
	} else {
		fileStatus += " [SAVED]"
	}

	modeStatus := " VIEW "
	if st.Mode == state.InsertMode {
		modeStatus = " INSERT "
	}

	leftStatus := modeStatus + fileStatus
	if message, ok := st.ActiveMessage(); ok {
		leftStatus = modeStatus + message
	}

	rightStatus := fmt.Sprintf("Row %d Col %d ", st.CurrentRow+1, st.CurrentCol)

	padding := st.Cols - len(leftStatus) - len(rightStatus)
	if padding < 0 {
		padding = 0
	}

	fullStatusBar := leftStatus + strings.Repeat(" ", padding) + rightStatus

	printMessage(0, st.Rows, theme.StatusBarForeground, theme.StatusBarBackground, fullStatusBar)
}

func printMessage(col, row int, foreground, background termbox.Attribute, message string) {
	for _, ch := range message {
		termbox.SetCell(col, row, ch, foreground, background)
		col += runewidth.RuneWidth(ch)
	}
}
