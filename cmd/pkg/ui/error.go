package ui

import (
	"strings"

	"github.com/nsf/termbox-go"
)

func ShowErrorMessage(message string) {
	width, _ := termbox.Size()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	lines := strings.Split(message, "\n")

	for y, line := range lines {
		startCol := (width - len(line)) / 2
		for i, ch := range line {
			termbox.SetCell(startCol+i, y, ch, termbox.ColorRed, termbox.ColorBlack)
		}
	}

	termbox.Flush()
}
