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
		if st.Mode == state.EditMode {
			st.Mode = state.ViewMode
			st.QuitKey = termbox.KeyEsc
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
