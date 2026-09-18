package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
)

func InsertRune(st *state.State, r rune) {
	if len(st.TextBuffer) == 0 {
		st.TextBuffer = append(st.TextBuffer, []rune{})
	}
	st.Do("insert", func() {
		line := st.TextBuffer[st.CurrentRow]
		newLine := append(line[:st.CurrentCol], append([]rune{r}, line[st.CurrentCol:]...)...)
		st.TextBuffer[st.CurrentRow] = newLine
		st.CurrentCol++
		st.Modified = true
	})
}

func InsertTab(st *state.State) {
	InsertRune(st, ' ')
}

func InsertNewLine(st *state.State) {
	if len(st.TextBuffer) == 0 {
		st.TextBuffer = append(st.TextBuffer, []rune{})
	}
	st.Do("newline", func() {
		line := st.TextBuffer[st.CurrentRow]
		beforeCursor := line[:st.CurrentCol]
		afterCursor := line[st.CurrentCol:]

		st.TextBuffer[st.CurrentRow] = beforeCursor
		st.TextBuffer = append(
			st.TextBuffer[:st.CurrentRow+1],
			append([][]rune{afterCursor}, st.TextBuffer[st.CurrentRow+1:]...)...,
		)
		st.CurrentRow++
		st.CurrentCol = 0
		st.Modified = true
	})
}

func Backspace(st *state.State) {
	if st.CurrentCol > 0 {
		st.Do("backspace", func() {
			line := st.TextBuffer[st.CurrentRow]
			newLine := append(line[:st.CurrentCol-1], line[st.CurrentCol:]...)
			st.TextBuffer[st.CurrentRow] = newLine
			st.CurrentCol--
			st.Modified = true
		})
	} else if st.CurrentRow > 0 {
		st.Do("backspace", func() {
			prevLine := st.TextBuffer[st.CurrentRow-1]
			currentLine := st.TextBuffer[st.CurrentRow]
			joinCol := len(prevLine)
			st.TextBuffer[st.CurrentRow-1] = append(prevLine, currentLine...)
			st.TextBuffer = append(st.TextBuffer[:st.CurrentRow], st.TextBuffer[st.CurrentRow+1:]...)
			st.CurrentRow--
			st.CurrentCol = joinCol
			st.Modified = true
		})
	}
}

func DeleteCurrentLine(st *state.State) {
	if st.CurrentRow >= len(st.TextBuffer) {
		return
	}

	st.Do("delete line", func() {
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

		st.Modified = true
	})
}
