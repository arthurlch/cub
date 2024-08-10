package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
)

// Adjusts the cursor column to the end of the line if it exceeds the line length
func adjustCursorColToLineEnd(s *state.State) {
	// Check if TextBuffer is empty
	if len(s.TextBuffer) == 0 {
		s.CurrentRow = 0
		s.CurrentCol = 0
		return
	}

	// Ensure CurrentRow is within bounds
	if s.CurrentRow >= len(s.TextBuffer) {
		s.CurrentRow = len(s.TextBuffer) - 1
	}

	// Ensure CurrentCol is within 
	if s.CurrentRow < len(s.TextBuffer) && s.CurrentCol > len(s.TextBuffer[s.CurrentRow]) {
		s.CurrentCol = len(s.TextBuffer[s.CurrentRow])
	}

	// Additional safeguard to ensure CurrentCol is never negative
	if s.CurrentCol < 0 {
		s.CurrentCol = 0
	}
}
