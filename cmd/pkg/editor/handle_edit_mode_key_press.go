package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

func handleEditModeKeyPress(es *EditorState, keyEvent termbox.Event) {
	st := es.State
	updateSelection(st)

	switch keyEvent.Key {
	case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight,
		termbox.KeyHome, termbox.KeyEnd, termbox.KeyPgup, termbox.KeyPgdn:
		handleNavigation(st, keyEvent)
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
	default:
		if keyEvent.Ch != 0 {
			if st.CurrentRow >= len(st.TextBuffer) {
				return
			}
			es.InsertRunes(keyEvent)
			st.Modified = true
		}
	}

	if st.CurrentRow < len(st.TextBuffer) && st.CurrentCol > len(st.TextBuffer[st.CurrentRow]) {
		st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
	}
	adjustCursorColToLineEnd(st)
	utils.ScrollTextBuffer(st)
}
