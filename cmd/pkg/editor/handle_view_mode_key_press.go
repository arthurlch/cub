package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

func handleViewModeKeyPress(st *state.State, keyEvent termbox.Event) {
	utils.LogKeyPress("handleViewModeKeyPress", keyEvent)

	switch keyEvent.Key {
	case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight,
		termbox.KeyHome, termbox.KeyEnd, termbox.KeyPgup, termbox.KeyPgdn:
		handleNavigation(st, keyEvent)
	case termbox.Key('z'):
		utils.Logger.Println("Starting selection")
		startSelection(st)
	case termbox.Key('a'):
		utils.Logger.Println("Selecting entire line")
		startSelection(st)
		st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
		updateSelection(st)
	case termbox.Key('s'):
		utils.Logger.Println("Ending selection")
		endSelection(st)
	case termbox.Key('c'):
		utils.Logger.Println("Copying selection")
		copySelection(st)
	case termbox.Key('x'):
		utils.Logger.Println("Cutting selection")
		cutSelection(st)
	case termbox.Key('v'):
		utils.Logger.Println("Pasting selection")
		pasteSelection(st)
	}
	adjustCursorColToLineEnd(st)
	utils.ScrollTextBuffer(st)
	utils.LogBufferState(st, "ViewMode")
}
