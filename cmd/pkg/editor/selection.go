package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
)

func startSelection(st *state.State, mode state.SelectionMode) {
	st.StartRow = st.CurrentRow
	st.StartCol = st.CurrentCol
	st.SelectionMode = mode
}

func updateSelection(st *state.State) {
	if st.SelectionMode == state.NoSelection {
		return
	}
}

func endSelection(st *state.State) {
	st.SelectionMode = state.NoSelection
}

func copySelection(st *state.State) {
	if st.SelectionMode != state.NoSelection {
		copyBuffer := []rune{}
		for row := st.StartRow; row <= st.CurrentRow; row++ {
			line := st.TextBuffer[row]
			if row == st.StartRow && row == st.CurrentRow {
				copyBuffer = append(copyBuffer, line[st.StartCol:st.CurrentCol]...)
			} else if row == st.StartRow {
				copyBuffer = append(copyBuffer, line[st.StartCol:]...)
			} else if row == st.CurrentRow {
				copyBuffer = append(copyBuffer, line[:st.CurrentCol]...)
			} else {
				copyBuffer = append(copyBuffer, line...)
			}
		}
		st.CopyBuffer = copyBuffer
		st.SelectionMode = state.NoSelection
	}
}

func cutSelection(st *state.State) {
	copySelection(st)
	if st.SelectionMode != state.NoSelection {
		newTextBuffer := [][]rune{}
		for row := 0; row < len(st.TextBuffer); row++ {
			if row < st.StartRow || row > st.CurrentRow {
				newTextBuffer = append(newTextBuffer, st.TextBuffer[row])
			} else {
				line := st.TextBuffer[row]
				if row == st.StartRow && row == st.CurrentRow {
					newTextBuffer = append(newTextBuffer, append(line[:st.StartCol], line[st.CurrentCol:]...))
				} else if row == st.StartRow {
					newTextBuffer = append(newTextBuffer, line[:st.StartCol])
				} else if row == st.CurrentRow {
					newTextBuffer = append(newTextBuffer, line[st.CurrentCol:])
				}
			}
		}
		st.TextBuffer = newTextBuffer
		st.Modified = true
		st.SelectionMode = state.NoSelection
	}
}

func pasteSelection(st *state.State) {
	if len(st.CopyBuffer) > 0 {
		line := st.TextBuffer[st.CurrentRow]
		before := line[:st.CurrentCol]
		after := line[st.CurrentCol:]
		newLine := append(append(before, st.CopyBuffer...), after...)
		st.TextBuffer[st.CurrentRow] = newLine
		st.Modified = true
	}
}
