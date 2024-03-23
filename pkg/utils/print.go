package utils

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func PrintMessage(col, row int, forground, background termbox.Attribute, message string) {
	for _, ch := range message {
		termbox.SetCell(col, row, ch, forground, background)
		col += runewidth.RuneWidth(ch)
	}
}