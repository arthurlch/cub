package editor

import (
	"os"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/ui"
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

func (es *EditorState) ProcessKeyPress() {
	st := es.State
	keyEvent := utils.GetKey()

	switch keyEvent.Type {
	case termbox.EventKey:
		handleKeyPress(es, keyEvent)
	case termbox.EventResize:
		st.Cols, st.Rows = termbox.Size()
		st.Rows--
		termbox.Flush()
	default:
	}
}

func handleKeyPress(es *EditorState, keyEvent termbox.Event) {
	st := es.State

	if keyEvent.Key == termbox.KeyEsc {
		if st.Mode == state.EditMode {
			st.Mode = state.ViewMode
			st.QuitKey = termbox.KeyEsc
			termbox.SetCursor(st.CurrentCol-st.OffsetCol, st.CurrentRow-st.OffsetRow)
			termbox.Flush()
		}
		return
	}

	if st.QuitKey == termbox.KeyEsc && keyEvent.Ch == 'q' {
		termbox.Close()
		os.Exit(0)
	} else {
		st.QuitKey = 0
	}

	if keyEvent.Ch == 'e' && st.Mode == state.ViewMode {
		st.Mode = state.EditMode
		return
	}

	if keyEvent.Key == termbox.KeyCtrlS {
		if err := es.SaveFile(); err != nil {
			ui.ShowErrorMessage(st, "Failed to save file: "+err.Error())
		} else {
			ui.ShowSuccessMessage(st, "File saved successfully.")
		}
		return
	}

	switch st.Mode {
	case state.ViewMode:
		handleViewModeKeyPress(st, keyEvent)
	case state.EditMode:
		handleEditModeKeyPress(es, keyEvent)
	}
}

func handleViewModeKeyPress(st *state.State, keyEvent termbox.Event) {
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
	}
	utils.ScrollTextBuffer(st)
}

func adjustCursorColToLineEnd(st *state.State) {
	if st.CurrentCol > len(st.TextBuffer[st.CurrentRow]) {
		st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
	}
}
