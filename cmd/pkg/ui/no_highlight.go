package ui

import (
	"github.com/nsf/termbox-go"
)

func RenderLineWithoutHighlight(line []rune, row int) {
	for col, ch := range line {
			termbox.SetCell(col, row, ch, TextForeground, TextBackground)
	}

	width, _ := termbox.Size()
	for col := len(line); col < width; col++ {
			termbox.SetCell(col, row, ' ', TextForeground, TextBackground)
	}
}

