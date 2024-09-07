package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

func handleInsertModeKeyPress(es *EditorState, keyEvent termbox.Event) {
	st := es.State
	updateSelection(st)

	if len(st.TextBuffer) == 0 {
		st.TextBuffer = append(st.TextBuffer, []rune{}) // add an empty line if the buffer is empty
	}

	if st.CurrentRow >= len(st.TextBuffer) {
		st.CurrentRow = len(st.TextBuffer) - 1
	}
	if st.CurrentRow < 0 {
		st.CurrentRow = 0
	}
	if st.CurrentCol > len(st.TextBuffer[st.CurrentRow]) {
		st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
	}
	if st.CurrentCol < 0 {
		st.CurrentCol = 0
	}

	switch keyEvent.Key {
	case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight,
		termbox.KeyHome, termbox.KeyEnd, termbox.KeyPgup, termbox.KeyPgdn:
		handleNavigation(st, keyEvent)
		utils.AdjustCursorColToLineEnd(st)
	case termbox.KeyEnter:
		currentLine := st.TextBuffer[st.CurrentRow]
		beforeCursor := currentLine[:st.CurrentCol]
		afterCursor := currentLine[st.CurrentCol:]

		st.TextBuffer[st.CurrentRow] = beforeCursor

		st.TextBuffer = append(st.TextBuffer[:st.CurrentRow+1], append([][]rune{afterCursor}, st.TextBuffer[st.CurrentRow+1:]...)...)

		st.CurrentRow++
		st.CurrentCol = 0

		recordChange(st, state.Change{Type: state.Insert, Row: st.CurrentRow, Col: 0, Text: []rune{'\n'}})
		st.Modified = true
	case termbox.KeyTab:
		recordChange(st, state.Change{Type: state.Insert, Row: st.CurrentRow, Col: st.CurrentCol, Text: []rune{'\t'}})
		for i := 0; i < 4; i++ {
			es.InsertRunes(keyEvent)
		}
		st.Modified = true
	case termbox.KeySpace:
		currentLine := st.TextBuffer[st.CurrentRow]
		st.TextBuffer[st.CurrentRow] = append(currentLine[:st.CurrentCol], append([]rune{' '}, currentLine[st.CurrentCol:]...)...)
		
		st.CurrentCol++

		recordChange(st, state.Change{Type: state.Insert, Row: st.CurrentRow, Col: st.CurrentCol - 1, Text: []rune{' '}})
		st.Modified = true
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		if st.CurrentRow < len(st.TextBuffer) && st.CurrentCol > 0 {
			deletedText := st.TextBuffer[st.CurrentRow][st.CurrentCol-1 : st.CurrentCol]
			recordChange(st, state.Change{Type: state.Delete, Row: st.CurrentRow, Col: st.CurrentCol - 1, Text: deletedText})
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
			recordChange(st, state.Change{Type: state.Insert, Row: st.CurrentRow, Col: st.CurrentCol, Text: []rune{keyEvent.Ch}})
			es.InsertRunes(keyEvent)
			st.Modified = true
		}
	}

	utils.AdjustCursorColToLineEnd(st)
	utils.ScrollTextBuffer(st)
}

func recordChange(s *state.State, change state.Change) {
	if s.HistoryIndex < len(s.ChangeHistory) {
		s.ChangeHistory = s.ChangeHistory[:s.HistoryIndex]
	}
	s.ChangeHistory = append(s.ChangeHistory, change)
	s.HistoryIndex++
}
