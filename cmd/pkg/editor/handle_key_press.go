package editor

import (
	"os"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/ui"
	"github.com/nsf/termbox-go"
)

func handleKeyPress(es *EditorState, keyEvent termbox.Event) {
	st := es.State

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


	if keyEvent.Ch == 'i' && st.Mode == state.ViewMode {
		st.Mode = state.InsertMode
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
