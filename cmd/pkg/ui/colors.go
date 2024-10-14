package ui

import (
	"github.com/nsf/termbox-go"
)

var (
	SoftBlack     = termbox.RGBToAttribute(30, 30, 30)     
	White         = termbox.RGBToAttribute(240, 240, 240)  
	LightGray     = termbox.RGBToAttribute(180, 180, 180)  
	Magenta       = termbox.RGBToAttribute(255, 85, 255)   
	Green         = termbox.RGBToAttribute(120, 255, 120)  
	Yellow        = termbox.RGBToAttribute(255, 215, 0)    
	Red           = termbox.RGBToAttribute(255, 70, 70)    
	Cyan          = termbox.RGBToAttribute(0, 255, 255)    
	Blue          = termbox.RGBToAttribute(70, 130, 180)   
	PinkBold      = termbox.RGBToAttribute(255, 105, 180)  
	ColorDarkPink = termbox.RGBToAttribute(231, 84, 128)  
	GoBlue        = termbox.RGBToAttribute(66, 165, 245)
)

var (
	ColorBackground     = SoftBlack
	ColorForeground     = White
	ColorHighlight      = Magenta
	StatusBarForeground = White
	StatusBarBackground = ColorDarkPink
	CursorForeground    = SoftBlack
	CursorBackground    = Green
	TextForeground      = White  
	TextBackground      = SoftBlack
	ModalTextColor      = GoBlue

	ColorKeyword  = PinkBold
	ColorString   = Green
	ColorComment  = LightGray
	ColorNumber   = Yellow
	ColorOperator = Red
	ColorDefault  = White
)
