package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/nsf/termbox-go"
)

func (es *EditorState) InsertRunes(event termbox.Event) {
	st := es.State
	if st == nil || st.CurrentRow >= len(st.TextBuffer) {
		return
	}

	row := st.TextBuffer[st.CurrentRow]
	if st.CurrentCol > len(row) {
		st.CurrentCol = len(row)
	}

	newRow := make([]rune, len(row)+1)
	copy(newRow, row[:st.CurrentCol])

	if event.Key == termbox.KeySpace {
		newRow[st.CurrentCol] = ' '
	} else {
		newRow[st.CurrentCol] = event.Ch
	}

	copy(newRow[st.CurrentCol+1:], row[st.CurrentCol:])
	st.TextBuffer[st.CurrentRow] = newRow
	st.CurrentCol++
}

func (es *EditorState) DeleteRune() {
	st := es.State
	if st == nil || st.CurrentRow >= len(st.TextBuffer) {
		return
	}

	row := st.TextBuffer[st.CurrentRow]
	if st.CurrentCol > 0 {
		st.CurrentCol--
		if len(row) > 0 {
			st.TextBuffer[st.CurrentRow] = append(row[:st.CurrentCol], row[st.CurrentCol+1:]...)
		}
	} else if st.CurrentRow > 0 {
		prevLineLen := len(st.TextBuffer[st.CurrentRow-1])
		st.TextBuffer[st.CurrentRow-1] = append(st.TextBuffer[st.CurrentRow-1], row...)
		st.TextBuffer = append(st.TextBuffer[:st.CurrentRow], st.TextBuffer[st.CurrentRow+1:]...)
		st.CurrentRow--
		st.CurrentCol = prevLineLen
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
	if st.CurrentRow < len(st.TextBuffer) {
		st.TextBuffer = append(st.TextBuffer[:st.CurrentRow], st.TextBuffer[st.CurrentRow+1:]...)
		if st.CurrentRow >= len(st.TextBuffer) {
			st.CurrentRow = len(st.TextBuffer) - 1
		}
		if st.CurrentRow < 0 {
			st.CurrentRow = 0
		}
		st.CurrentCol = 0
		st.Modified = true
	}
}
