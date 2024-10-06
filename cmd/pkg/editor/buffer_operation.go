package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
)

func Undo(s *state.State) {
	if len(s.UndoBuffer) == 0 {
			return
	}

	s.RedoBuffer = append(s.RedoBuffer, state.UndoState{
			TextBuffer: utils.DeepCopyTextBuffer(s.TextBuffer),
			CurrentRow: s.CurrentRow,
			CurrentCol: s.CurrentCol,
	})

	lastIndex := len(s.UndoBuffer) - 1
	lastState := s.UndoBuffer[lastIndex]
	s.TextBuffer = lastState.TextBuffer
	s.CurrentRow = lastState.CurrentRow
	s.CurrentCol = lastState.CurrentCol
	s.UndoBuffer = s.UndoBuffer[:lastIndex]

	utils.ValidateCursorPosition(s)
	s.Modified = true
}

func Redo(s *state.State) {
	if len(s.RedoBuffer) == 0 {
			return
	}

	s.UndoBuffer = append(s.UndoBuffer, state.UndoState{
			TextBuffer: utils.DeepCopyTextBuffer(s.TextBuffer),
			CurrentRow: s.CurrentRow,
			CurrentCol: s.CurrentCol,
	})

	lastIndex := len(s.RedoBuffer) - 1
	lastState := s.RedoBuffer[lastIndex]
	s.TextBuffer = lastState.TextBuffer
	s.CurrentRow = lastState.CurrentRow
	s.CurrentCol = lastState.CurrentCol
	s.RedoBuffer = s.RedoBuffer[:lastIndex]

	utils.ValidateCursorPosition(s)
	s.Modified = true
}
