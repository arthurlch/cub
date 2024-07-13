package utils

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
)

func ScrollTextBuffer(s *state.State) {
	if s.CurrentRow < s.OffsetRow {
		s.OffsetRow = s.CurrentRow
	}

	if s.CurrentRow >= s.OffsetRow+s.Rows {
		s.OffsetRow = s.CurrentRow - s.Rows + 1
	}

	if s.CurrentCol < s.OffsetCol {
		s.OffsetCol = s.CurrentCol
	}

	if s.CurrentCol >= s.OffsetCol+s.Cols {
		s.OffsetCol = s.CurrentCol - s.Cols + 1
		if s.OffsetCol < 0 {
			s.OffsetCol = 0
		}
	}
}
