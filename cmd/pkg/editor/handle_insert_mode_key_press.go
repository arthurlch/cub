package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

func handleInsertModeKeyPress(es *EditorState, keyEvent termbox.Event) {
	st := es.State
	updateSelection(st)

	switch keyEvent.Key {
	case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight,
		termbox.KeyHome, termbox.KeyEnd, termbox.KeyPgup, termbox.KeyPgdn:
		handleNavigation(st, keyEvent)
		utils.AdjustCursorColToLineEnd(st)
	case termbox.KeyEnter:
		saveChangeToUndoBuffer(st, state.Change{Type: state.Insert, Row: st.CurrentRow, Col: st.CurrentCol, Text: []rune{'\n'}})
		es.InsertNewLine()
		st.Modified = true
	case termbox.KeyTab:
		saveChangeToUndoBuffer(st, state.Change{Type: state.Insert, Row: st.CurrentRow, Col: st.CurrentCol, Text: []rune{'\t'}})
		for i := 0; i < 4; i++ {
			es.InsertRunes(keyEvent)
		}
		st.Modified = true
	case termbox.KeySpace:
		saveChangeToUndoBuffer(st, state.Change{Type: state.Insert, Row: st.CurrentRow, Col: st.CurrentCol, Text: []rune{' '}})
		es.InsertRunes(keyEvent)
		st.Modified = true
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		if st.CurrentRow < len(st.TextBuffer) && st.CurrentCol > 0 {
			deletedText := st.TextBuffer[st.CurrentRow][st.CurrentCol-1 : st.CurrentCol]
			saveChangeToUndoBuffer(st, state.Change{Type: state.Delete, Row: st.CurrentRow, Col: st.CurrentCol - 1, Text: deletedText})
			es.DeleteRune()
			st.Modified = true
		} else if st.CurrentCol == 0 && st.CurrentRow > 0 {
			prevRowLength := len(st.TextBuffer[st.CurrentRow-1])
			st.CurrentCol = prevRowLength
			st.TextBuffer[st.CurrentRow-1] = append(st.TextBuffer[st.CurrentRow-1], st.TextBuffer[st.CurrentRow]...)
			st.TextBuffer = append(st.TextBuffer[:st.CurrentRow], st.TextBuffer[st.CurrentRow+1:]...)
			st.CurrentRow--
			st.Modified = true
		}
	default:
		if keyEvent.Ch != 0 {
			saveChangeToUndoBuffer(st, state.Change{Type: state.Insert, Row: st.CurrentRow, Col: st.CurrentCol, Text: []rune{keyEvent.Ch}})
			es.InsertRunes(keyEvent)
			st.Modified = true
		}
	}

	utils.AdjustCursorColToLineEnd(st)
	utils.ScrollTextBuffer(st)
}
