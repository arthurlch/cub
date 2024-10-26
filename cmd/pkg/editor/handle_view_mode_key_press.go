package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

func handleViewModeKeyPress(st *state.State, keyEvent termbox.Event) {
	if st.SelectionActive && (keyEvent.Key == termbox.KeyDelete || keyEvent.Key == termbox.KeyBackspace || keyEvent.Key == termbox.KeyBackspace2) {
		DeleteSelection(st)
		return
	}

	switch keyEvent.Ch {
	case 's':
		StartSelection(st)
	case 'z':
		EndSelection(st)
	case 'c':
		CopySelection(st)
		EndSelection(st)
	case 'x':
		CutSelection(st)
		EndSelection(st)
	case 'v':
		PasteSelection(st)
	case 'd':
		if st.LastKey == 'd' {
			DeleteCurrentLine(st)
			st.LastKey = 0
		} else {
			st.LastKey = 'd'
		}
	case 'a':
		if st.LastKey == 'a' {
			SelectAll(st) 
			st.LastKey = 0
		} else {
			st.LastKey = 'a'
		}

	default:
		handleSpecialKeys(st, keyEvent)
		handleComplexNavigation(st, keyEvent)
	}

	if st.SelectionActive && isNavigationKey(keyEvent) {
		UpdateSelection(st)
	}

	utils.AdjustCursorColToLineEnd(st)
	utils.ScrollTextBuffer(st)
}

func handleSpecialKeys(st *state.State, keyEvent termbox.Event) {
	switch keyEvent.Key {
	case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight,
		termbox.KeyHome, termbox.KeyEnd, termbox.KeyPgup, termbox.KeyPgdn:
		handleNavigation(st, keyEvent)
	default:
		st.LastKey = 0
	}
}

func isNavigationKey(keyEvent termbox.Event) bool {
	switch keyEvent.Key {
	case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight,
		termbox.KeyHome, termbox.KeyEnd, termbox.KeyPgup, termbox.KeyPgdn:
		return true
	}
	return false
}
