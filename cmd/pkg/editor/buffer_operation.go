package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
)

func Undo(s *state.State) {
	if s.HistoryIndex <= 0 {
		return
	}
	s.HistoryIndex--
	change := s.ChangeHistory[s.HistoryIndex]
	inverseChange := invertChange(change)
	utils.ApplyChange(s, inverseChange)

	s.CurrentRow = change.Row
	s.CurrentCol = change.Col

	utils.AdjustCursorColToLineEnd(s)
}

func Redo(s *state.State) {
	if s.HistoryIndex >= len(s.ChangeHistory) {
		return
	}
	change := s.ChangeHistory[s.HistoryIndex]
	s.HistoryIndex++
	utils.ApplyChange(s, change)

	s.CurrentRow = change.Row
	s.CurrentCol = change.Col

	utils.AdjustCursorColToLineEnd(s)
}

func invertChange(c state.Change) state.Change {
	return state.Change{
		Type:    invertChangeType(c.Type),
		Row:     c.Row,
		Col:     c.Col,
		Text:    c.Text,
		PrevRow: c.PrevRow,
		PrevCol: c.PrevCol,
	}
}

func invertChangeType(t state.ChangeType) state.ChangeType {
	if t == state.Insert {
		return state.Delete
	}
	return state.Insert
}