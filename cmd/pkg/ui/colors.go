package ui

import "github.com/nsf/termbox-go"

const (
	ColorBackground   = termbox.ColorDefault 
	ColorForeground   = termbox.ColorMagenta 
	ColorHighlight    = termbox.ColorWhite
)

var (
	StatusBarForeground = ColorForeground
	StatusBarBackground = termbox.ColorMagenta  
	CursorForeground    = ColorBackground
	CursorBackground    = termbox.ColorRed  
	TextForeground      = ColorForeground
	TextBackground      = ColorBackground
)
