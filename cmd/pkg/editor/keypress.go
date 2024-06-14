package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

func handleEditModeKeyPress(es *EditorState, keyEvent termbox.Event) {
	st := es.State
	updateSelection(st)
	switch keyEvent.Key {
	case termbox.KeyArrowUp:
		if st.CurrentRow > 0 {
			st.CurrentRow--
			adjustCursorColToLineEnd(st)
		}
	case termbox.KeyArrowDown:
		if st.CurrentRow < len(st.TextBuffer)-1 {
			st.CurrentRow++
			adjustCursorColToLineEnd(st)
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
		adjustCursorColToLineEnd(st)
	case termbox.KeyPgdn:
		if st.CurrentRow+int(st.Rows/4) < len(st.TextBuffer)-1 {
			st.CurrentRow += int(st.Rows / 4)
		} else {
			st.CurrentRow = len(st.TextBuffer) - 1
		}
		adjustCursorColToLineEnd(st)
	case termbox.KeyEnter:
		es.InsertNewLine()
		st.Modified = true
	case termbox.KeyTab:
		for i := 0; i < 4; i++ {
			es.InsertRunes(keyEvent)
		}
		st.Modified = true
	case termbox.KeySpace:
		es.InsertRunes(keyEvent)
		st.Modified = true
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		es.DeleteRune()
		st.Modified = true
	}

	if keyEvent.Ch != 0 {
		if st.CurrentRow >= len(st.TextBuffer) {
			return
		}
		es.InsertRunes(keyEvent)
		st.Modified = true
	}

	if st.CurrentRow < len(st.TextBuffer) && st.CurrentCol > len(st.TextBuffer[st.CurrentRow]) {
		st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
	}
	adjustCursorColToLineEnd(st)
	utils.ScrollTextBuffer(st)
}

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
