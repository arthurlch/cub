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
		st.TextBuffer = append(st.TextBuffer, []rune{})
	}

	utils.ValidateCursorPosition(st)

	switch keyEvent.Key {
	case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight,
		termbox.KeyHome, termbox.KeyEnd, termbox.KeyPgup, termbox.KeyPgdn:
		handleNavigation(st, keyEvent)

	case termbox.KeyEnter:
		saveUndoState(st)
		currentLine := st.TextBuffer[st.CurrentRow]
		beforeCursor := currentLine[:st.CurrentCol]
		afterCursor := currentLine[st.CurrentCol:]

		st.TextBuffer[st.CurrentRow] = beforeCursor
		st.TextBuffer = append(
			st.TextBuffer[:st.CurrentRow+1],
			append([][]rune{afterCursor}, st.TextBuffer[st.CurrentRow+1:]...)...,
		)

		st.CurrentRow++
		st.CurrentCol = 0
		st.Modified = true

	case termbox.KeyTab:
		saveUndoState(st)
		tabSpaces := []rune{' '}
		line := st.TextBuffer[st.CurrentRow]
		newLine := append(line[:st.CurrentCol], append(tabSpaces, line[st.CurrentCol:]...)...)
		st.TextBuffer[st.CurrentRow] = newLine
		st.CurrentCol += len(tabSpaces)
		st.Modified = true

	case termbox.KeySpace:
		saveUndoState(st)
		insertRune(st, ' ')

	case termbox.KeyBackspace, termbox.KeyBackspace2:
		handleBackspace(st)

	default:
		if keyEvent.Ch != 0 {
			saveUndoState(st)
			insertRune(st, keyEvent.Ch)
		}
	}

	utils.ValidateCursorPosition(st) 
	utils.ScrollTextBuffer(st)       
}

func saveUndoState(st *state.State) {
	st.UndoBuffer = append(st.UndoBuffer, state.UndoState{
		TextBuffer: utils.DeepCopyTextBuffer(st.TextBuffer),
		CurrentRow: st.CurrentRow,
		CurrentCol: st.CurrentCol,
	})
	st.RedoBuffer = nil
}

func insertRune(st *state.State, r rune) {
	line := st.TextBuffer[st.CurrentRow]
	newLine := append(line[:st.CurrentCol], append([]rune{r}, line[st.CurrentCol:]...)...)
	st.TextBuffer[st.CurrentRow] = newLine
	st.CurrentCol++
	st.Modified = true
}

func handleBackspace(st *state.State) {
	if st.CurrentCol > 0 {
		saveUndoState(st)
		line := st.TextBuffer[st.CurrentRow]
		newLine := append(line[:st.CurrentCol-1], line[st.CurrentCol:]...)
		st.TextBuffer[st.CurrentRow] = newLine
		st.CurrentCol--
		st.Modified = true
	} else if st.CurrentCol == 0 && st.CurrentRow > 0 {
		saveUndoState(st)
		prevLine := st.TextBuffer[st.CurrentRow-1]
		currentLine := st.TextBuffer[st.CurrentRow]
		st.TextBuffer[st.CurrentRow-1] = append(prevLine, currentLine...)
		st.TextBuffer = append(st.TextBuffer[:st.CurrentRow], st.TextBuffer[st.CurrentRow+1:]...)
		st.CurrentRow--
		st.CurrentCol = len(prevLine)
		st.Modified = true
	}
}
