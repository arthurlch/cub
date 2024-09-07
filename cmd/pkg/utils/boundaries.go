package utils

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
)

func ApplyChange(s *state.State, c state.Change) {
	EnsurePositionExists(s, c.Row, c.Col)
	switch c.Type {
	case state.Insert:
		for i, ch := range c.Text {
			s.TextBuffer[c.Row] = append(s.TextBuffer[c.Row][:c.Col+i], append([]rune{ch}, s.TextBuffer[c.Row][c.Col+i:]...)...)
		}
	case state.Delete:
		if c.Row < len(s.TextBuffer) && c.Col+len(c.Text) <= len(s.TextBuffer[c.Row]) {
			s.TextBuffer[c.Row] = append(s.TextBuffer[c.Row][:c.Col], s.TextBuffer[c.Row][c.Col+len(c.Text):]...)
		}
	}
	s.CurrentRow, s.CurrentCol = c.PrevRow, c.PrevCol
}

func EnsurePositionExists(s *state.State, row, col int) {
	for len(s.TextBuffer) <= row {
		s.TextBuffer = append(s.TextBuffer, []rune{})
	}
	for len(s.TextBuffer[row]) <= col {
		s.TextBuffer[row] = append(s.TextBuffer[row], ' ')
	}
}

func EnsureRowExists(s *state.State, row int) {
	for len(s.TextBuffer) <= row {
		s.TextBuffer = append(s.TextBuffer, []rune{})
	}
}

func EnsureColExists(s *state.State, row, col int) {
	EnsureRowExists(s, row)
	for len(s.TextBuffer[row]) <= col {
		s.TextBuffer[row] = append(s.TextBuffer[row], ' ')
	}
}

func AdjustCursorColToLineEnd(s *state.State) {
	if s.CurrentRow < len(s.TextBuffer) && s.CurrentCol > len(s.TextBuffer[s.CurrentRow]) {
		s.CurrentCol = len(s.TextBuffer[s.CurrentRow])
	}
}
