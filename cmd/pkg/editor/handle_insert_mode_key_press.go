package editor

import (
	"github.com/nsf/termbox-go"
)

func handleInsertModeKeyPress(es *EditorState, keyEvent termbox.Event) {
	st := es.State

	UpdateSelection(st)

	if len(st.TextBuffer) == 0 {
		st.TextBuffer = append(st.TextBuffer, []rune{})
	}

	st.ClampCursor()

	switch keyEvent.Key {
	case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight,
		termbox.KeyHome, termbox.KeyEnd, termbox.KeyPgup, termbox.KeyPgdn:
		handleNavigation(st, keyEvent)

	case termbox.KeyEnter:
		InsertNewLine(st)

	case termbox.KeyTab:
		InsertTab(st)

	case termbox.KeySpace:
		InsertRune(st, ' ')

	case termbox.KeyBackspace, termbox.KeyBackspace2:
		Backspace(st)

	default:
		if keyEvent.Ch != 0 {
			InsertRune(st, keyEvent.Ch)
		}
	}

	st.ClampCursor()
	st.Scroll()
}
