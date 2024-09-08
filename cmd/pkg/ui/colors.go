package ui

import "github.com/nsf/termbox-go"

const (
	ColorBackground   = termbox.ColorBlack
	ColorForeground   = termbox.ColorWhite
	ColorHighlight    = termbox.ColorMagenta
	ColorPurple       = termbox.ColorMagenta
	ColorGreen        = termbox.ColorGreen
	ColorOrange       = termbox.ColorYellow
	ColorRed          = termbox.ColorRed
)

var (
	StatusBarForeground = ColorBackground
	StatusBarBackground = ColorPurple
	CursorForeground    = ColorBackground
	CursorBackground    = ColorOrange
	TextForeground      = ColorForeground
	TextBackground      = ColorBackground
)