package utils

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/nsf/termbox-go"
)

func DisplayTextBuffer(s *state.State) {
	var row, col int
	for row = 0; row < s.Rows; row++ {
			textBufferRow := row + s.OffsetRow
			for col = 0; col < s.Cols; col++ {
					textBufferCol := col + s.OffsetCol
					if textBufferRow >= 0 && textBufferRow < len(s.TextBuffer) && textBufferCol < len(s.TextBuffer[textBufferRow]) {
							char := s.TextBuffer[textBufferRow][textBufferCol]
							fgColor := termbox.ColorDefault
							bgColor := termbox.ColorDefault

							if s.SelectionActive && isSelected(s, textBufferRow, textBufferCol) {
									fgColor = termbox.ColorBlack
									bgColor = termbox.ColorWhite
							}

							if char != '\t' {
									termbox.SetCell(col, row, char, fgColor, bgColor)
							} else {
									termbox.SetCell(col, row, rune(' '), fgColor, bgColor)
							}
					} else if row+s.OffsetCol > len(s.TextBuffer) {
							termbox.SetCell(0, row, rune('*'), termbox.ColorLightMagenta, termbox.ColorDefault)
							termbox.SetChar(col, row, rune('\n'))
					}
			}
	}

	termbox.SetCursor(s.CurrentCol-s.OffsetCol, s.CurrentRow-s.OffsetRow)
}

func isSelected(st *state.State, row, col int) bool {
	startRow, endRow := st.StartRow, st.EndRow
	startCol, endCol := st.StartCol, st.EndCol

	if startRow > endRow || (startRow == endRow && startCol > endCol) {
			startRow, endRow = endRow, startRow
			startCol, endCol = endCol, startCol
	}

	if row < startRow || row > endRow {
			return false
	}
	if row == startRow && col < startCol {
			return false
	}
	if row == endRow && col >= endCol {
			return false
	}
	return true
}
