package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

func handleViewModeKeyPress(st *state.State, keyEvent termbox.Event) {
	utils.LogKeyPress("handleViewModeKeyPress", keyEvent)
	utils.Logger.Printf("Current state before key press: SelectionActive=%v, Row=%d, Col=%d", 
		st.SelectionActive, st.CurrentRow, st.CurrentCol)

	switch {
	case keyEvent.Key == termbox.KeyArrowUp || keyEvent.Ch == 'o':
		handleNavigation(st, termbox.Event{Key: termbox.KeyArrowUp})
	case keyEvent.Key == termbox.KeyArrowDown || keyEvent.Ch == 'p':
		handleNavigation(st, termbox.Event{Key: termbox.KeyArrowDown})
	case keyEvent.Key == termbox.KeyArrowLeft || keyEvent.Ch == 'k':
		handleNavigation(st, termbox.Event{Key: termbox.KeyArrowLeft})
	case keyEvent.Key == termbox.KeyArrowRight || keyEvent.Ch == 'l':
		handleNavigation(st, termbox.Event{Key: termbox.KeyArrowRight})
	case keyEvent.Key == termbox.KeyHome || keyEvent.Key == termbox.KeyEnd || 
	     keyEvent.Key == termbox.KeyPgup || keyEvent.Key == termbox.KeyPgdn:
		handleNavigation(st, keyEvent)
	case keyEvent.Ch == 'z':
		utils.Logger.Printf("Attempting to start selection")
		startSelection(st)
	case keyEvent.Ch == 's':
		utils.Logger.Printf("Attempting to end selection")
		endSelection(st)
	case keyEvent.Ch == 'c':
		utils.Logger.Printf("Attempting to copy selection")
		copySelection(st)
		endSelection(st)
	case keyEvent.Ch == 'x':
		utils.Logger.Printf("Attempting to cut selection")
		cutSelection(st)
		endSelection(st)
	case keyEvent.Ch == 'v':
		utils.Logger.Printf("Attempting to paste selection")
		pasteSelection(st)
	case keyEvent.Ch == 'd':
		if st.LastKey == 'd' {
			utils.Logger.Printf("Attempting to delete current line")
			deleteCurrentLine(st)
			st.LastKey = 0
		} else {
			st.LastKey = 'd'
		}
	default:
		utils.Logger.Printf("Unhandled key press in view mode: %v", keyEvent)
		st.LastKey = 0
	}

	if st.SelectionActive && 
	   (keyEvent.Key == termbox.KeyArrowUp || keyEvent.Key == termbox.KeyArrowDown || 
	    keyEvent.Key == termbox.KeyArrowLeft || keyEvent.Key == termbox.KeyArrowRight || 
	    keyEvent.Key == termbox.KeyHome || keyEvent.Key == termbox.KeyEnd || 
	    keyEvent.Key == termbox.KeyPgup || keyEvent.Key == termbox.KeyPgdn) {
		utils.Logger.Printf("Updating selection during navigation")
		updateSelection(st)
	}

	utils.Logger.Printf("Current state after key press: SelectionActive=%v, Row=%d, Col=%d", 
		st.SelectionActive, st.CurrentRow, st.CurrentCol)
	
	adjustCursorColToLineEnd(st)
	utils.ScrollTextBuffer(st)
	utils.LogBufferState(st, "ViewMode")
}
