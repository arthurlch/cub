package syntax

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/arthurlch/cub/cmd/pkg/ui"
	"github.com/nsf/termbox-go"
)

func GetLexer(fileType string) chroma.Lexer {
    lexer := lexers.Get(fileType)
    if lexer == nil {
        return lexers.Fallback
    }
    return chroma.Coalesce(lexer)
}

func GetTermboxColor(tokenType chroma.TokenType, tokenValue string) (termbox.Attribute, termbox.Attribute) {
    bg := ui.TextBackground

    switch tokenType.Category() {
    case chroma.Keyword:
        return ui.Blue, bg
    case chroma.Name:
        if tokenType == chroma.NameFunction {
            return ui.Yellow, bg
        }
        return ui.White, bg
    case chroma.String:
        return ui.PinkBold, bg
    case chroma.Number:
        return ui.Yellow, bg
    case chroma.Comment:
        return ui.LightGray, bg
    case chroma.Operator, chroma.Punctuation:
        return ui.Red, bg
    default:
        return ui.White, bg
    }
}