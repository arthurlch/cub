package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
)

func Undo(s *state.State) {
	if s.Undo() {
		s.Modified = true
	}
}

func Redo(s *state.State) {
	if s.Redo() {
		s.Modified = true
	}
}
