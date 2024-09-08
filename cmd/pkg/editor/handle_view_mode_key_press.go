package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

func handleViewModeKeyPress(st *state.State, keyEvent termbox.Event) {

	if st.SelectionActive && (keyEvent.Key == termbox.KeyDelete || keyEvent.Key == termbox.KeyBackspace || keyEvent.Key == termbox.KeyBackspace2) {
		deleteSelection(st) 
		return
	}

	switch keyEvent.Ch {
	case 'o', 'p', 'k', 'l':
		handleNavigationAlias(st, keyEvent.Ch)
		if st.SelectionActive {
			updateSelection(st)  
		}
	case 's':
		startSelection(st)
	case 'z':
		endSelection(st)
	case 'c':
		copySelection(st)
		endSelection(st)
	case 'x':
		cutSelection(st)
		endSelection(st)
	case 'v':
		pasteSelection(st)
	case 'd':
		if st.LastKey == 'd' {
			deleteCurrentLine(st)
			st.LastKey = 0
		} else {
			st.LastKey = 'd'
		}
	default:
		handleSpecialKeys(st, keyEvent)
	}

	if st.SelectionActive && isNavigationKey(keyEvent) {
		updateSelection(st)
	}

	adjustCursorColToLineEnd(st)
	utils.ScrollTextBuffer(st)

	
}

/* 

Decision making for this one was tough but ergonomically speaking 
Having that position is better for the wrist, problem with fully horinzontal setup
is that the wrist is getting tired quickly. 

*/ 

func handleNavigationAlias(st *state.State, ch rune) {
	switch ch {
	case 'o':
		handleNavigation(st, termbox.Event{Key: termbox.KeyArrowUp})
	case 'p':
		handleNavigation(st, termbox.Event{Key: termbox.KeyArrowDown})
	case 'k':
		handleNavigation(st, termbox.Event{Key: termbox.KeyArrowLeft})
	case 'l':
		handleNavigation(st, termbox.Event{Key: termbox.KeyArrowRight})
	}
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
