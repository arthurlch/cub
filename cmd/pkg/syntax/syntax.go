package syntax

import (
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/ui"
	"github.com/nsf/termbox-go"
)

func RenderHighlightedLine(line string, row int, fileType string, sharedState *state.State) {
	lexer := lexers.Get(fileType)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	iterator, _ := lexer.Tokenise(nil, line)
	col := 0

	for token := iterator(); token != chroma.EOF; token = iterator() {
		fg, bg := getTermboxColor(token.Type, token.Value)
		for _, ch := range token.Value {
			if col < sharedState.Cols {
				termbox.SetCell(col, row, ch, fg, bg)
				col++
			}
		}
	}
}

func getTermboxColor(tokenType chroma.TokenType, tokenValue string) (termbox.Attribute, termbox.Attribute) {
	if isOperatorOrPunctuation(tokenValue) {
		return RGBToAttribute(198, 120, 221), termbox.ColorDefault
	}

	switch {
	case tokenType.InCategory(chroma.Keyword):
		return RGBToAttribute(255, 121, 198), termbox.ColorDefault

	case tokenType == chroma.NameFunction:
		return RGBToAttribute(241, 250, 140), termbox.ColorDefault

	case tokenType.InCategory(chroma.String):
		return RGBToAttribute(195, 232, 141), termbox.ColorDefault 

	case tokenType.InCategory(chroma.Comment):
		return RGBToAttribute(98, 114, 164), termbox.ColorDefault 

	case tokenType.InCategory(chroma.LiteralNumber):
		return RGBToAttribute(247, 140, 108), termbox.ColorDefault 

	default:
		return RGBToAttribute(248, 248, 242), termbox.ColorDefault
	}
}

func isOperatorOrPunctuation(value string) bool {
	operators := []string{
		"=", "==", "+", "-", "*", "/", "%", "<", ">", "!=", "<=", ">=", "&&", "||", ":",
	}
	for _, operator := range operators {
		if value == operator {
			return true
		}
	}
	return false
}

func IsSupportedLanguage(fileType string) bool {
	return lexers.Get(fileType) != nil
}

func RGBToAttribute(r, g, b uint8) termbox.Attribute {
	return termbox.RGBToAttribute(r, g, b)
}

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
