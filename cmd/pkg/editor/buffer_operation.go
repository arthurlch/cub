package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
)


var maxUndoStates = 100
// SaveChangeToUndoBuffer saves a change to the undo buffer with cursor position
func saveChangeToUndoBuffer(s *state.State, change state.Change) {
	change.PrevRow = s.CurrentRow
	change.PrevCol = s.CurrentCol
	s.UndoBuffer = append(s.UndoBuffer, change)
	if len(s.UndoBuffer) > maxUndoStates {
		s.UndoBuffer = s.UndoBuffer[1:]
	}
	s.RedoBuffer = nil // Clear the redo buffer on new change
}

// Undo reverts the last change made to the text buffer
func Undo(s *state.State) {
	utils.Logger.Println("Undo called")
	if len(s.UndoBuffer) == 0 {
		utils.Logger.Println("UndoBuffer is empty")
		return
	}

	lastChange := s.UndoBuffer[len(s.UndoBuffer)-1]
	s.UndoBuffer = s.UndoBuffer[:len(s.UndoBuffer)-1]

	inverseChange := invertChange(lastChange)
	s.RedoBuffer = append(s.RedoBuffer, lastChange)
	utils.ApplyChange(s, inverseChange)
	utils.Logger.Printf("UndoBuffer length: %d", len(s.UndoBuffer))
	utils.Logger.Printf("RedoBuffer length: %d", len(s.RedoBuffer))
	utils.LogTextBuffer(s.TextBuffer, "After Undo")
}

// Redo reapplies the last undone change to the text buffer
func Redo(s *state.State) {
	utils.Logger.Println("Redo called")
	if len(s.RedoBuffer) == 0 {
		utils.Logger.Println("RedoBuffer is empty")
		return
	}

	lastChange := s.RedoBuffer[len(s.RedoBuffer)-1]
	s.RedoBuffer = s.RedoBuffer[:len(s.RedoBuffer)-1]

	s.UndoBuffer = append(s.UndoBuffer, lastChange)
	utils.ApplyChange(s, lastChange)
	utils.Logger.Printf("UndoBuffer length: %d", len(s.UndoBuffer))
	utils.Logger.Printf("RedoBuffer length: %d", len(s.RedoBuffer))
	utils.LogTextBuffer(s.TextBuffer, "After Redo")
}

// invertChange creates the inverse of a given change
func invertChange(c state.Change) state.Change {
	inverseChange := state.Change{
		Type:    c.Type,
		Row:     c.Row,
		Col:     c.Col,
		Text:    c.Text,
		PrevRow: c.PrevRow,
		PrevCol: c.PrevCol,
	}
	switch c.Type {
	case state.Insert:
		inverseChange.Type = state.Delete
	case state.Delete:
		inverseChange.Type = state.Insert
	}
	return inverseChange
}
