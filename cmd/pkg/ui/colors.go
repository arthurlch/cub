package ui

import (
	"github.com/arthurlch/cub/cmd/pkg/syntax"
	"github.com/nsf/termbox-go"
)

const (
    ColorBackground   = termbox.ColorBlack
    ColorForeground   = termbox.ColorWhite
    ColorHighlight    = termbox.ColorMagenta
    ColorPurple       = termbox.ColorMagenta
    ColorGreen        = termbox.ColorGreen
    ColorOrange       = termbox.ColorYellow
    ColorRed          = termbox.ColorRed
    ColorCyan         = termbox.ColorCyan
		ColorBlack        = termbox.ColorBlack
    ColorBlue         = termbox.ColorBlue
    ColorPink         = termbox.ColorMagenta | termbox.AttrBold
)

var (
    StatusBarForeground = termbox.ColorWhite
    StatusBarBackground = ColorPurple
    CursorForeground    = ColorBlack
    CursorBackground    = ColorGreen
    TextForeground      = ColorForeground
    TextBackground      = ColorBackground
    ColorKeyword        = ColorPink   
    ColorString         = ColorGreen  
    ColorComment        = ColorCyan   
    ColorNumber         = ColorOrange 
    ColorOperator       = ColorRed    
    ColorDefault        = ColorForeground
)

func HighlightAndRenderLine(line string, y int) {
	tokens := syntax.Tokenize(line)
	x := 0
	for _, token := range tokens {
		var fg termbox.Attribute
		switch token.Type {
		case syntax.TokenKeyword:
				fg = ColorKeyword
		case syntax.TokenString:
				fg = ColorString
		case syntax.TokenComment:
				fg = ColorComment
		case syntax.TokenNumber:
				fg = ColorNumber
		case syntax.TokenOperator:
				fg = ColorOperator
		case syntax.TokenWhitespace:
				fg = ColorDefault
		default:
				fg = ColorDefault
		}

		for _, ch := range token.Value {
				termbox.SetCell(x, y, ch, fg, ColorBackground)
				x++
		}
	}
}
