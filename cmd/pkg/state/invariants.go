package state

func (s *State) ClampCursor() {
	if len(s.TextBuffer) == 0 {
		s.CurrentRow, s.CurrentCol = 0, 0
		return
	}

	if s.CurrentRow < 0 {
		s.CurrentRow = 0
	} else if s.CurrentRow >= len(s.TextBuffer) {
		s.CurrentRow = len(s.TextBuffer) - 1
	}

	if s.CurrentCol < 0 {
		s.CurrentCol = 0
	} else if s.CurrentCol > len(s.TextBuffer[s.CurrentRow]) {
		s.CurrentCol = len(s.TextBuffer[s.CurrentRow])
	}
}

func (s *State) MoveCursor(newRow, newCol int) {
	if len(s.TextBuffer) == 0 {
		s.CurrentRow, s.CurrentCol = 0, 0
		return
	}

	if newRow < 0 {
		s.CurrentRow = 0
	} else if newRow >= len(s.TextBuffer) {
		s.CurrentRow = len(s.TextBuffer) - 1
	} else {
		s.CurrentRow = newRow
	}

	lineLength := len(s.TextBuffer[s.CurrentRow])
	if newCol < 0 {
		s.CurrentCol = 0
	} else if newCol > lineLength {
		s.CurrentCol = lineLength
	} else {
		s.CurrentCol = newCol
	}
}

func (s *State) SnapCursorToLineEnd() {
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

func (s *State) EnsurePositionExists(row, col int) {
	s.EnsureRowExists(row)
	s.EnsureColExists(row, col)
}

func (s *State) EnsureRowExists(row int) {
	for len(s.TextBuffer) <= row {
		s.TextBuffer = append(s.TextBuffer, []rune{})
	}
}

func (s *State) EnsureColExists(row, col int) {
	s.EnsureRowExists(row)
	for len(s.TextBuffer[row]) <= col {
		s.TextBuffer[row] = append(s.TextBuffer[row], ' ')
	}
}

func (s *State) Scroll() {
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

func DeepCopyTextBuffer(buffer [][]rune) [][]rune {
	newBuffer := make([][]rune, len(buffer))
	for i, row := range buffer {
		newRow := make([]rune, len(row))
		copy(newRow, row)
		newBuffer[i] = newRow
	}
	return newBuffer
}
