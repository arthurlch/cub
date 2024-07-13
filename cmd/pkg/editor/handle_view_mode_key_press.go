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
		if st.SelectionActive {
			utils.Logger.Printf("Updating selection during navigation")
			updateSelection(st)
		}
	case keyEvent.Key == termbox.KeyArrowDown || keyEvent.Ch == 'p':
		handleNavigation(st, termbox.Event{Key: termbox.KeyArrowDown})
		if st.SelectionActive {
			utils.Logger.Printf("Updating selection during navigation")
			updateSelection(st)
		}
	case keyEvent.Key == termbox.KeyArrowLeft || keyEvent.Ch == 'k':
		handleNavigation(st, termbox.Event{Key: termbox.KeyArrowLeft})
		if st.SelectionActive {
			utils.Logger.Printf("Updating selection during navigation")
			updateSelection(st)
		}
	case keyEvent.Key == termbox.KeyArrowRight || keyEvent.Ch == 'l':
		handleNavigation(st, termbox.Event{Key: termbox.KeyArrowRight})
		if st.SelectionActive {
			utils.Logger.Printf("Updating selection during navigation")
			updateSelection(st)
		}
	case keyEvent.Key == termbox.KeyHome || keyEvent.Key == termbox.KeyEnd || 
	     keyEvent.Key == termbox.KeyPgup || keyEvent.Key == termbox.KeyPgdn:
		handleNavigation(st, keyEvent)
		if st.SelectionActive {
			utils.Logger.Printf("Updating selection during navigation")
			updateSelection(st)
		}
	case keyEvent.Ch == 'z':
		utils.Logger.Printf("Attempting to start selection")
		startSelection(st)
		utils.Logger.Printf("Selection started: Active=%v, StartRow=%d, StartCol=%d", 
			st.SelectionActive, st.StartRow, st.StartCol)
	case keyEvent.Ch == 's':
		utils.Logger.Printf("Attempting to end selection")
		endSelection(st)
		utils.Logger.Printf("Selection ended: Active=%v", st.SelectionActive)
	case keyEvent.Ch == 'c':
		utils.Logger.Printf("Attempting to copy selection")
		copySelection(st)
		utils.Logger.Printf("Copy complete, buffer length: %d", len(st.CopyBuffer))
		endSelection(st)
	case keyEvent.Ch == 'x':
		utils.Logger.Printf("Attempting to cut selection")
		cutSelection(st)
		utils.Logger.Printf("Cut complete, buffer length: %d", len(st.CopyBuffer))
		endSelection(st)

	case keyEvent.Ch == 'v':
		utils.Logger.Printf("Attempting to paste selection")
		pasteSelection(st)
		utils.Logger.Printf("Paste complete")

	case keyEvent.Ch == 'd':
		if st.LastKey == 'd' {
			utils.Logger.Printf("Attempting to delete current line")
			deleteCurrentLine(st)
			utils.Logger.Printf("Current line deleted")
			st.LastKey = 0
		} else {
			st.LastKey = 'd'
		}

	default:
		utils.Logger.Printf("Unhandled key press in view mode: %v", keyEvent)
		st.LastKey = 0
	}

	utils.Logger.Printf("Current state after key press: SelectionActive=%v, Row=%d, Col=%d", 
		st.SelectionActive, st.CurrentRow, st.CurrentCol)
	
	adjustCursorColToLineEnd(st)
	utils.ScrollTextBuffer(st)
	utils.LogBufferState(st, "ViewMode")
}
