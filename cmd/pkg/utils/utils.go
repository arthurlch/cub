package utils

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/nsf/termbox-go"
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

func DisplayTextBuffer(s *state.State) {
	var row, col int
	for row = 0; row < s.Rows; row++ {
		textBufferRow := row + s.OffsetRow
		for col = 0; col < s.Cols; col++ {
			textBufferCol := col + s.OffsetCol
			if textBufferRow >= 0 && textBufferRow < len(s.TextBuffer) && textBufferCol < len(s.TextBuffer[textBufferRow]) {
				if s.TextBuffer[textBufferRow][textBufferCol] != '\t' {
					termbox.SetChar(col, row, s.TextBuffer[textBufferRow][textBufferCol])
				} else {
					termbox.SetCell(col, row, rune(' '), termbox.ColorDefault, termbox.ColorDefault)
				}
			} else if row+s.OffsetCol > len(s.TextBuffer) {
				termbox.SetCell(0, row, rune('*'), termbox.ColorLightMagenta, termbox.ColorDefault)
				termbox.SetChar(col, row, rune('\n'))
			}
		}
	}
}

func GetKey() termbox.Event {
	var key_event termbox.Event
	switch event := termbox.PollEvent(); event.Type {
	case termbox.EventKey:
		key_event = event
	case termbox.EventError:
		panic(event.Err)
	}
	return key_event
}
