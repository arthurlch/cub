package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/logging"
	"github.com/arthurlch/cub/cmd/pkg/state"
)

func StartSelection(st *state.State) {
	st.StartRow = st.CurrentRow
	st.StartCol = st.CurrentCol
	st.EndRow = st.CurrentRow
	st.EndCol = st.CurrentCol
	st.SelectionActive = true
	logging.Logger.Printf("Selection started - StartRow: %d, StartCol: %d, SelectionActive: %v",
		st.StartRow, st.StartCol, st.SelectionActive)
}

func UpdateSelection(st *state.State) {
	st.EndRow = st.CurrentRow
	st.EndCol = st.CurrentCol

	st.ClampCursor()

	if st.StartRow == st.EndRow && st.StartCol <= st.EndCol {
		st.EndCol++
	}
}

func EndSelection(st *state.State) {
	st.SelectionActive = false
	logging.Logger.Println("End selection")
}

func CopySelection(st *state.State) {
	copyBuffer := []rune{}
	startRow, endRow := st.StartRow, st.EndRow
	startCol, endCol := st.StartCol, st.EndCol

	if startRow > endRow || (startRow == endRow && startCol > endCol) {
		startRow, endRow = endRow, startRow
		startCol, endCol = endCol, startCol
	}

	if startCol > len(st.TextBuffer[startRow]) {
		startCol = len(st.TextBuffer[startRow])
	}
	if endCol > len(st.TextBuffer[endRow]) {
		endCol = len(st.TextBuffer[endRow])
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
	st.CopyBuffer = copyBuffer
	logging.Logger.Printf("Copy selection - CopyBuffer length: %d", len(st.CopyBuffer))
}

func CutSelection(st *state.State) {
	logging.Logger.Println("Cut selection - Start")
	CopySelection(st)
	DeleteSelection(st)
	logging.Logger.Printf("Cut selection - Removed text length: %d", len(st.CopyBuffer))
}

func PasteSelection(st *state.State) {
	if len(st.CopyBuffer) == 0 {
		logging.Logger.Println("Paste selection - No text to paste")
		return
	}

	st.Do("paste", func() {
		lines := [][]rune{{}}
		for _, ch := range st.CopyBuffer {
			if ch == '\n' {
				lines = append(lines, []rune{})
			} else {
				lines[len(lines)-1] = append(lines[len(lines)-1], ch)
			}
		}

		currentLine := st.TextBuffer[st.CurrentRow]
		before := append([]rune{}, currentLine[:st.CurrentCol]...)
		after := append([]rune{}, currentLine[st.CurrentCol:]...)

		newTextBuffer := append([][]rune{}, st.TextBuffer[:st.CurrentRow]...)
		newTextBuffer = append(newTextBuffer, append(before, lines[0]...))

		for i := 1; i < len(lines); i++ {
			newTextBuffer = append(newTextBuffer, lines[i])
		}

		newTextBuffer[len(newTextBuffer)-1] = append(newTextBuffer[len(newTextBuffer)-1], after...)

		newTextBuffer = append(newTextBuffer, st.TextBuffer[st.CurrentRow+1:]...)

		st.TextBuffer = newTextBuffer
		st.CurrentRow += len(lines) - 1
		st.CurrentCol = len(newTextBuffer[st.CurrentRow]) - len(after)
		st.Modified = true

		st.ClampCursor()
	})
}

func DeleteSelection(st *state.State) {
	st.Do("delete selection", func() {
		newTextBuffer := [][]rune{}
		startRow, endRow := st.StartRow, st.EndRow
		startCol, endCol := st.StartCol, st.EndCol

		if startRow > endRow || (startRow == endRow && startCol > endCol) {
			startRow, endRow = endRow, startRow
			startCol, endCol = endCol, startCol
		}

		for row := 0; row < len(st.TextBuffer); row++ {
			if row < startRow || row > endRow {
				newTextBuffer = append(newTextBuffer, st.TextBuffer[row])
			} else if row == startRow && row == endRow {
				st.EnsurePositionExists(row, endCol)
				line := st.TextBuffer[row]
				newLine := append(line[:startCol], line[endCol:]...)
				newTextBuffer = append(newTextBuffer, newLine)
			} else if row == startRow {
				st.EnsurePositionExists(row, startCol)
				line := st.TextBuffer[row]
				newTextBuffer = append(newTextBuffer, line[:startCol])
			} else if row == endRow {
				st.EnsurePositionExists(row, endCol)
				line := st.TextBuffer[row]
				newTextBuffer = append(newTextBuffer, line[endCol:])
			}
		}

		if len(newTextBuffer) == 0 {
			newTextBuffer = append(newTextBuffer, []rune{})
		}
		st.TextBuffer = newTextBuffer

		st.CurrentRow = startRow
		st.CurrentCol = startCol

		st.ClampCursor()

		st.Modified = true
	})
}

// missing select all
func SelectAll(st *state.State) {

	if len(st.TextBuffer) == 0 {
		return
	}

	st.StartRow = 0
	st.StartCol = 0
	st.EndRow = len(st.TextBuffer) - 1
	st.EndCol = len(st.TextBuffer[st.StartRow])

	st.SelectionActive = true
}
