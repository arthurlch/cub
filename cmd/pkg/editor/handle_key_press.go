package editor

import (
	"os"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/ui"
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

func handleKeyPress(es *EditorState, keyEvent termbox.Event) {
	st := es.State


	if keyEvent.Key == termbox.KeyCtrlH {
		ui.ShowHelpModal()
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
		termbox.Close()
		os.Exit(0)
	}

	if keyEvent.Key == termbox.KeyCtrlU {
		utils.Logger.Println("Ctrl+U pressed")
		Undo(st)
		return
	}

	if keyEvent.Key == termbox.KeyCtrlR {
		utils.Logger.Println("Ctrl+R pressed")
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
			ui.ShowErrorMessage(st, "Failed to save file: "+err.Error())
		} else {
			ui.ShowSuccessMessage(st, "File saved successfully.")
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
