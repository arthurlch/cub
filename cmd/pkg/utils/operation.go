package utils

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
)

func CollectSelectionText(st *state.State) []rune {
	copyBuffer := []rune{}
	startRow, endRow := st.StartRow, st.EndRow
	startCol, endCol := st.StartCol, st.EndCol

	if startRow > endRow || (startRow == endRow && startCol > endCol) {
		startRow, endRow = endRow, startRow
		startCol, endCol = endCol, startCol
	}

	for row := startRow; row <= endRow; row++ {
		line := st.TextBuffer[row]
		if row == startRow && row == endRow {
			copyBuffer = append(copyBuffer, line[startCol:endCol]...)
		} else if row == startRow {
			copyBuffer = append(copyBuffer, line[startCol:]...)
			copyBuffer = append(copyBuffer, '\n')
		} else if row == endRow {
			copyBuffer = append(copyBuffer, line[:endCol]...)
		} else {
			copyBuffer = append(copyBuffer, line...)
			copyBuffer = append(copyBuffer, '\n')
		}
	}
	return copyBuffer
}
