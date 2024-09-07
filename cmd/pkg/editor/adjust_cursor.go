package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
)

func adjustCursorColToLineEnd(s *state.State) {
	if len(s.TextBuffer) == 0 {
		s.CurrentRow = 0
		s.CurrentCol = 0
		return
	}

	if s.CurrentRow < 0 {
		s.CurrentRow = 0
	}

	if s.CurrentRow >= len(s.TextBuffer) {
		s.CurrentRow = len(s.TextBuffer) - 1
	}

	if s.CurrentCol > len(s.TextBuffer[s.CurrentRow]) {
		s.CurrentCol = len(s.TextBuffer[s.CurrentRow])
	}

	if s.CurrentCol < 0 {
		s.CurrentCol = 0
	}
}
