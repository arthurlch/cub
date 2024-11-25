package syntax

import (
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
	"github.com/arthurlch/cub/cmd/pkg/ui"
	"github.com/nsf/termbox-go"
)

func GetLexer(fileType string) chroma.Lexer {
	lexer := lexers.Get(fileType)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	return lexer
}

func GetTermboxColor(tokenType chroma.TokenType, tokenValue string) (termbox.Attribute, termbox.Attribute) {
	fg := ui.ColorDefault
	bg := ui.ColorBackground

	if isOperatorOrPunctuation(tokenValue) {
		fg = ui.ColorOperator
	} else {
		switch {
		case tokenType.InCategory(chroma.Keyword):
			fg = ui.ColorKeyword
		case tokenType == chroma.NameFunction:
			fg = ui.ColorNumber 
		case tokenType.InCategory(chroma.String):
			fg = ui.ColorString
		case tokenType.InCategory(chroma.Comment):
			fg = ui.ColorComment
		case tokenType.InCategory(chroma.LiteralNumber):
			fg = ui.ColorNumber
		default:
			fg = ui.ColorDefault
		}
	}
	return fg, bg
}

func isOperatorOrPunctuation(value string) bool {
	operators := map[string]bool{
		"=": true, "==": true, "+": true, "-": true, "*": true,
		"/": true, "%": true, "<": true, ">": true, "!=": true,
		"<=": true, ">=": true, "&&": true, "||": true, ":": true,
	}
	return operators[value]
}