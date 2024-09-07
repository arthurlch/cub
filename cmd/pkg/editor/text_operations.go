package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/nsf/termbox-go"
)

func (es *EditorState) InsertRunes(keyEvent termbox.Event) {
	st := es.State

	if len(st.TextBuffer) == 0 {
		st.TextBuffer = append(st.TextBuffer, []rune{}) 
	}

	if st.CurrentRow < 0 || st.CurrentRow >= len(st.TextBuffer) {
		st.CurrentRow = 0 
	}

	if st.CurrentCol < 0 {
		st.CurrentCol = 0 
	}

	if st.CurrentCol > len(st.TextBuffer[st.CurrentRow]) {
		st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
	}

	if keyEvent.Ch != 0 {
		st.TextBuffer[st.CurrentRow] = append(
			st.TextBuffer[st.CurrentRow][:st.CurrentCol],
			append([]rune{keyEvent.Ch}, st.TextBuffer[st.CurrentRow][st.CurrentCol:]...)...,
		)
		st.CurrentCol++
	}
}

func (es *EditorState) DeleteRune() {
	st := es.State
	if st.CurrentRow < len(st.TextBuffer) && st.CurrentCol > 0 {
		row := st.TextBuffer[st.CurrentRow]
		st.TextBuffer[st.CurrentRow] = append(row[:st.CurrentCol-1], row[st.CurrentCol:]...)
		st.CurrentCol--
	} else if st.CurrentCol == 0 && st.CurrentRow > 0 {
		prevRowLength := len(st.TextBuffer[st.CurrentRow-1])
		st.CurrentCol = prevRowLength
		st.TextBuffer[st.CurrentRow-1] = append(st.TextBuffer[st.CurrentRow-1], st.TextBuffer[st.CurrentRow]...)
		st.TextBuffer = append(st.TextBuffer[:st.CurrentRow], st.TextBuffer[st.CurrentRow+1:]...)
		st.CurrentRow--
	}
}

func (es *EditorState) InsertNewLine() {
	st := es.State
	if st == nil || st.CurrentRow >= len(st.TextBuffer) {
		return
	}
	beforeCursor := st.TextBuffer[st.CurrentRow][:st.CurrentCol]
	afterCursor := st.TextBuffer[st.CurrentRow][st.CurrentCol:]

	st.TextBuffer[st.CurrentRow] = beforeCursor

	if st.CurrentRow+1 < len(st.TextBuffer) {
		st.TextBuffer = append(st.TextBuffer[:st.CurrentRow+1], append([][]rune{afterCursor}, st.TextBuffer[st.CurrentRow+1:]...)...)
	} else {
		st.TextBuffer = append(st.TextBuffer, afterCursor)
	}

	st.CurrentRow++
	st.CurrentCol = 0
}

func deleteCurrentLine(st *state.State) {
	if st.CurrentRow >= len(st.TextBuffer) {
		return
	}

	deletedLine := st.TextBuffer[st.CurrentRow]

	st.ChangeHistory = append(st.ChangeHistory[:st.HistoryIndex], state.Change{
		Type:    state.Delete,
		Row:     st.CurrentRow,
		Col:     0,
		Text:    deletedLine,
		PrevRow: st.CurrentRow,
		PrevCol: st.CurrentCol,
	})
	st.HistoryIndex++

	st.TextBuffer = append(st.TextBuffer[:st.CurrentRow], st.TextBuffer[st.CurrentRow+1:]...)

	if len(st.TextBuffer) == 0 {
		st.TextBuffer = append(st.TextBuffer, []rune{}) 
		st.CurrentRow = 0
		st.CurrentCol = 0
	} else {
		if st.CurrentRow >= len(st.TextBuffer) {
			st.CurrentRow = len(st.TextBuffer) - 1
		}
		if st.CurrentCol > len(st.TextBuffer[st.CurrentRow]) {
			st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
		}
	}
}
