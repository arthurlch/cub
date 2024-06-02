package editor

import (
	"log"
	"os"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/ui"
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

func (es *EditorState) ProcessKeyPress() {
	st := es.State
	log.Println("Waiting for key event...")
	keyEvent := utils.GetKey()
	log.Printf("Key event received: %+v\n", keyEvent)

	switch keyEvent.Type {
	case termbox.EventKey:
		handleKeyPress(es, keyEvent)
	case termbox.EventResize:
		// Handle terminal resize events if necessary
		st.Cols, st.Rows = termbox.Size()
		st.Rows--
		termbox.Flush()
	}
}

func handleKeyPress(es *EditorState, keyEvent termbox.Event) {
	st := es.State

	if keyEvent.Key == termbox.KeyEsc {
		if st.Mode == state.EditMode {
			st.Mode = state.ViewMode
			st.QuitKey = termbox.KeyEsc
			st.StopBlink <- struct{}{}
			termbox.SetCursor(st.CurrentCol-st.OffsetCol, st.CurrentRow-st.OffsetRow)
			termbox.Flush()
			log.Println("Switched to view mode.")
		}
		return
	}

	if st.QuitKey == termbox.KeyEsc && keyEvent.Ch == 'q' {
		close(st.StopBlink)
		termbox.Close()
		log.Println("Quitting editor.")
		os.Exit(0)
	} else {
		st.QuitKey = 0
	}

	if keyEvent.Ch == 'e' && st.Mode == state.ViewMode {
		st.Mode = state.EditMode
		go es.blinkCursor()
		log.Println("Switched to edit mode.")
		return
	}

	if keyEvent.Key == termbox.KeyCtrlS {
		log.Println("Saving file...")
		if err := es.SaveFile(st.SourceFile); err != nil {
			ui.ShowErrorMessage("Failed to save file: " + err.Error())
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
	log.Println("Handling view mode key press...")
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
	log.Println("View mode key press handled.")
}

func handleEditModeKeyPress(es *EditorState, keyEvent termbox.Event) {
	st := es.State
	log.Println("Handling edit mode key press...")
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
	case termbox.KeyEnter:
		es.InsertNewLine()
		st.Modified = true
	case termbox.KeyTab:
		for i := 0; i < 4; i++ {
			es.InsertRunes(keyEvent)
		}
		st.Modified = true
	case termbox.KeySpace:
		es.InsertRunes(keyEvent)
		st.Modified = true
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		es.DeleteRune()
		st.Modified = true
	}

	if keyEvent.Ch != 0 {
		if st.CurrentRow >= len(st.TextBuffer) {
			return
		}
		es.InsertRunes(keyEvent)
		st.Modified = true
	}

	if st.CurrentRow < len(st.TextBuffer) && st.CurrentCol > len(st.TextBuffer[st.CurrentRow]) {
		st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
	}
	utils.ScrollTextBuffer(st)
	log.Println("Edit mode key press handled.")
}

func adjustCursorColToLineEnd(st *state.State) {
	if st.CurrentCol > len(st.TextBuffer[st.CurrentRow]) {
		st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
	}
}
