package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/nsf/termbox-go"
)

func handleViewModeKeyPress(st *state.State, keyEvent termbox.Event) {
	if st.SelectionActive && isDeleteKey(keyEvent.Key) {
		DeleteSelection(st)
		return
	}

	switch {
	case tryViewKeyBinding(st, keyEvent):
	case tryPendingSequence(st, keyEvent):
	case tryViewRuneBinding(st, keyEvent):
	default:
		st.LastKey = 0
	}

	if st.SelectionActive && isNavigationKey(keyEvent) {
		UpdateSelection(st)
	}

	st.SnapCursorToLineEnd()
	st.Scroll()
}

func isNavigationKey(keyEvent termbox.Event) bool {
	switch keyEvent.Key {
	case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight,
		termbox.KeyHome, termbox.KeyEnd, termbox.KeyPgup, termbox.KeyPgdn:
		return true
	}
	return false
}
