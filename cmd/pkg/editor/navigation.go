package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

func handleNavigation(st *state.State, keyEvent termbox.Event) {
	switch keyEvent.Key {
	case termbox.KeyArrowUp:
		if st.CurrentRow > 0 {
			st.CurrentRow--
		}
	case termbox.KeyArrowDown:
		if st.CurrentRow < len(st.TextBuffer)-1 {
			st.CurrentRow++
		}
	case termbox.KeyArrowLeft:
		if st.CurrentCol > 0 {
			st.CurrentCol--
		} else if st.CurrentRow > 0 {
			st.CurrentRow--
			st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
		}
	case termbox.KeyArrowRight:
		if st.CurrentCol < len(st.TextBuffer[st.CurrentRow]) {
			st.CurrentCol++
		} else if st.CurrentRow < len(st.TextBuffer)-1 {
			st.CurrentRow++
			st.CurrentCol = 0
		}
	case termbox.KeyHome:
		st.CurrentCol = 0
	case termbox.KeyEnd:
		st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
	case termbox.KeyPgup:
		if st.CurrentRow-int(st.Rows/4) > 0 {
			st.CurrentRow -= int(st.Rows / 4)
		} else {
			st.CurrentRow = 0
		}
	case termbox.KeyPgdn:
		if st.CurrentRow+int(st.Rows/4) < len(st.TextBuffer)-1 {
			st.CurrentRow += int(st.Rows / 4)
		} else {
			st.CurrentRow = len(st.TextBuffer) - 1
		}
	}
	utils.LogKeyPress("handleNavigation", keyEvent)
	utils.LogBufferState(st, "Navigation")
}
