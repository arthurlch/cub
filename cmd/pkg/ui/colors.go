package ui

import (
	"github.com/nsf/termbox-go"
)

var (
    SoftBlack     = termbox.RGBToAttribute(1, 22, 39)      // #011627 Background
    White         = termbox.RGBToAttribute(214, 222, 235)   // #D6DEEB Foreground
    LightGray     = termbox.RGBToAttribute(99, 119, 119)    // #637777 Comment
    Magenta       = termbox.RGBToAttribute(199, 146, 234)   // #C792EA Keywords
    Green         = termbox.RGBToAttribute(173, 219, 103)   // #ADDB67 Support
    Yellow        = termbox.RGBToAttribute(247, 140, 108)   // #F78C6C Numbers
    Red           = termbox.RGBToAttribute(255, 99, 99)     // #FF6363 Constant
    Cyan          = termbox.RGBToAttribute(127, 219, 202)   // #7FDBCA Language variables
    Blue          = termbox.RGBToAttribute(130, 170, 255)   // #82AAFF Functions
    PinkBold      = termbox.RGBToAttribute(236, 196, 141)   // #ECC48D Quoted strings
    ColorDarkPink = termbox.RGBToAttribute(29, 59, 83)      // #1D3B53 Active selection
    GoBlue        = termbox.RGBToAttribute(128, 164, 194)   // #80A4C2 Cursor
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

    SelectedBackground = termbox.ColorWhite
    SelectedForeground = termbox.ColorWhite

    ColorKeyword  = PinkBold
    ColorString   = Green
    ColorComment  = LightGray
    ColorNumber   = Yellow
    ColorOperator = Red
    ColorDefault  = White
)