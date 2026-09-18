package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/logging"
	"github.com/arthurlch/cub/cmd/pkg/render"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/nsf/termbox-go"
)

func handleKeyPress(es *EditorState, keyEvent termbox.Event) {
	st := es.State

	if keyEvent.Key == termbox.KeyCtrlH {
		render.ShowHelpModal()
		return
	}

	if keyEvent.Key == termbox.KeyEsc {
		if st.Mode == state.InsertMode {
			st.Mode = state.ViewMode
			st.QuitKey = termbox.KeyEsc
		}
		return
	}

	if keyEvent.Key == termbox.KeyCtrlQ {
		st.Quit = true
		return
	}

	if keyEvent.Key == termbox.KeyCtrlU {
		logging.Logger.Println("Ctrl+U pressed")
		Undo(st)
		return
	}

	if keyEvent.Key == termbox.KeyCtrlR {
		logging.Logger.Println("Ctrl+R pressed")
		Redo(st)
		return
	}

	if keyEvent.Ch == 'i' && st.Mode == state.ViewMode {
		st.Mode = state.InsertMode
		EndSelection(st)
		if st.CurrentRow >= len(st.TextBuffer) {
			st.CurrentRow = len(st.TextBuffer) - 1
		}
		if st.CurrentRow >= 0 && st.CurrentCol > len(st.TextBuffer[st.CurrentRow]) {
			st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
		}
		return
	}

	if keyEvent.Key == termbox.KeyCtrlS {
		if err := es.SaveFile(); err != nil {
			st.ShowMessage("Failed to save file: " + err.Error())
		} else {
			st.ShowMessage("File saved successfully.")
		}
		return
	}

	switch st.Mode {
	case state.ViewMode:
		handleViewModeKeyPress(st, keyEvent)
	case state.InsertMode:
		handleInsertModeKeyPress(es, keyEvent)
	}
}
