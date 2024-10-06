package utils

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
)

func EnsurePositionExists(s *state.State, row, col int) {
    EnsureRowExists(s, row)
    EnsureColExists(s, row, col)
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

func AdjustCursorAfterChange(s *state.State, newRow, newCol int) {
	if newRow < 0 {
			s.CurrentRow = 0
	} else if newRow >= len(s.TextBuffer) {
			s.CurrentRow = len(s.TextBuffer) - 1
	} else {
			s.CurrentRow = newRow
	}

	lineLength := len(s.TextBuffer[s.CurrentRow])
	if newCol < 0 {
			s.CurrentCol = 0
	} else if newCol > lineLength {
			s.CurrentCol = lineLength
	} else {
			s.CurrentCol = newCol
	}
}

func ValidateCursorPosition(s *state.State) {
    if s.CurrentRow < 0 {
        s.CurrentRow = 0
    } else if s.CurrentRow >= len(s.TextBuffer) {
        s.CurrentRow = len(s.TextBuffer) - 1
    }

    if s.CurrentCol < 0 {
        s.CurrentCol = 0
    } else if s.CurrentCol > len(s.TextBuffer[s.CurrentRow]) {
        s.CurrentCol = len(s.TextBuffer[s.CurrentRow])
    }
}


func DeepCopyTextBuffer(buffer [][]rune) [][]rune {
	newBuffer := make([][]rune, len(buffer))
	for i, row := range buffer {
			newRow := make([]rune, len(row))
			copy(newRow, row)
			newBuffer[i] = newRow
	}
	return newBuffer
}

func AdjustCursorColToLineEnd(s *state.State) {
	if len(s.TextBuffer) == 0 {
		s.CurrentRow = 0
		s.CurrentCol = 0
		return
	}

	if s.CurrentRow < 0 {
		s.CurrentRow = 0
	}

	if s.CurrentRow >= len(s.TextBuffer) {
		s.CurrentRow = len(s.TextBuffer) - 1
	}

	if s.CurrentCol > len(s.TextBuffer[s.CurrentRow]) {
		s.CurrentCol = len(s.TextBuffer[s.CurrentRow])
	}

	if s.CurrentCol < 0 {
		s.CurrentCol = 0
	}
}
