package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
)

func startSelection(st *state.State) {
	st.StartRow = st.CurrentRow
	st.StartCol = st.CurrentCol
	utils.Logger.Printf("Start selection - StartRow: %d, StartCol: %d\n", st.StartRow, st.StartCol)
}

func updateSelection(st *state.State) {
	utils.Logger.Printf("Update selection - CurrentRow: %d, CurrentCol: %d\n", st.CurrentRow, st.CurrentCol)
}

func endSelection(st *state.State) {
	utils.Logger.Println("End selection")
}

func copySelection(st *state.State) {
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
	utils.Logger.Printf("Copy selection - CopyBuffer: %s\n", string(copyBuffer))
}

func cutSelection(st *state.State) {
	utils.Logger.Println("Cut selection - Start")
	copySelection(st)
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
	utils.Logger.Println("Cut selection - TextBuffer modified")
}

func pasteSelection(st *state.State) {
	if len(st.CopyBuffer) > 0 {
		line := st.TextBuffer[st.CurrentRow]
		before := line[:st.CurrentCol]
		after := line[st.CurrentCol:]
		newLine := append(append(before, st.CopyBuffer...), after...)
		st.TextBuffer[st.CurrentRow] = newLine
		st.Modified = true
		utils.Logger.Printf("Paste selection - Pasted text: %s\n", string(st.CopyBuffer))
	} else {
		utils.Logger.Println("Paste selection - No text to paste")
	}
}
