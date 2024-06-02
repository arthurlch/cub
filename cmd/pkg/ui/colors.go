package ui

import "github.com/nsf/termbox-go"

// Define the default color scheme using termbox colors by default
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
