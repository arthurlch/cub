package ui

import (
	"github.com/arthurlch/cub/cmd/pkg/theme"
	"github.com/nsf/termbox-go"
)

func RenderLineWithoutHighlight(line []rune, row int) {
	for col, ch := range line {
			termbox.SetCell(col, row, ch, theme.TextForeground, theme.TextBackground)
	}

	width, _ := termbox.Size()
	for col := len(line); col < width; col++ {
			termbox.SetCell(col, row, ' ', theme.TextForeground, theme.TextBackground)
	}
}

