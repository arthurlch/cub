package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
)

func startSelection(st *state.State) {
	st.StartRow = st.CurrentRow
	st.StartCol = st.CurrentCol
	st.EndRow = st.CurrentRow
	st.EndCol = st.CurrentCol
	st.SelectionActive = true
	utils.Logger.Printf("Selection started - StartRow: %d, StartCol: %d, SelectionActive: %v",
		st.StartRow, st.StartCol, st.SelectionActive)
}

func updateSelection(st *state.State) {
	st.EndRow = st.CurrentRow
	st.EndCol = st.CurrentCol

	utils.ValidateCursorPosition(st)

	if st.StartRow == st.EndRow && st.StartCol <= st.EndCol {
		st.EndCol++
	}
}

func endSelection(st *state.State) {
	st.SelectionActive = false
	utils.Logger.Println("End selection")
}

func copySelection(st *state.State) {
	copyBuffer := []rune{}
	startRow, endRow := st.StartRow, st.EndRow
	startCol, endCol := st.StartCol, st.EndCol

	if startRow > endRow || (startRow == endRow && startCol > endCol) {
		startRow, endRow = endRow, startRow
		startCol, endCol = endCol, startCol
	}

	utils.EnsurePositionExists(st, endRow, endCol)

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
	utils.Logger.Printf("Copy selection - CopyBuffer length: %d", len(st.CopyBuffer))
}

func cutSelection(st *state.State) {
	utils.Logger.Println("Cut selection - Start")
	copySelection(st)
	deleteSelection(st)
	utils.Logger.Printf("Cut selection - Removed text length: %d", len(st.CopyBuffer))
}

func pasteSelection(st *state.State) {
	if len(st.CopyBuffer) > 0 {
		st.UndoBuffer = append(st.UndoBuffer, state.UndoState{
			TextBuffer: utils.DeepCopyTextBuffer(st.TextBuffer),
			CurrentRow: st.CurrentRow,
			CurrentCol: st.CurrentCol,
		})
		st.RedoBuffer = nil

		lines := [][]rune{{}}
		for _, ch := range st.CopyBuffer {
			if ch == '\n' {
				lines = append(lines, []rune{})
			} else {
				lines[len(lines)-1] = append(lines[len(lines)-1], ch)
			}
		}

		currentLine := st.TextBuffer[st.CurrentRow]
		before := currentLine[:st.CurrentCol]
		after := currentLine[st.CurrentCol:]

		newTextBuffer := append([][]rune{}, st.TextBuffer[:st.CurrentRow]...)
		newTextBuffer = append(newTextBuffer, append(before, lines[0]...))

		for i := 1; i < len(lines); i++ {
			newTextBuffer = append(newTextBuffer, lines[i])
		}

		if len(lines) > 1 {
			newTextBuffer[len(newTextBuffer)-1] = append(newTextBuffer[len(newTextBuffer)-1], after...)
		} else {
			newTextBuffer[len(newTextBuffer)-1] = append(newTextBuffer[len(newTextBuffer)-1], after...)
		}

		newTextBuffer = append(newTextBuffer, st.TextBuffer[st.CurrentRow+1:]...)

		st.TextBuffer = newTextBuffer
		st.CurrentRow += len(lines) - 1
		st.CurrentCol = len(newTextBuffer[st.CurrentRow]) - len(after)
		st.Modified = true

		utils.ValidateCursorPosition(st)

	} else {
		utils.Logger.Println("Paste selection - No text to paste")
	}
}

func deleteSelection(st *state.State) {
	st.UndoBuffer = append(st.UndoBuffer, state.UndoState{
		TextBuffer: utils.DeepCopyTextBuffer(st.TextBuffer),
		CurrentRow: st.CurrentRow,
		CurrentCol: st.CurrentCol,
	})
	st.RedoBuffer = nil 

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
			utils.EnsurePositionExists(st, row, endCol)
			line := st.TextBuffer[row]
			newLine := append(line[:startCol], line[endCol:]...)
			newTextBuffer = append(newTextBuffer, newLine)
		} else if row == startRow {
			utils.EnsurePositionExists(st, row, startCol)
			line := st.TextBuffer[row]
			newTextBuffer = append(newTextBuffer, line[:startCol])
		} else if row == endRow {
			utils.EnsurePositionExists(st, row, endCol)
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

	utils.ValidateCursorPosition(st)

	st.Modified = true
}
