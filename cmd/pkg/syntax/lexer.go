package syntax

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/arthurlch/cub/cmd/pkg/theme"
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
    bg := theme.TextBackground

    switch tokenType.Category() {
    case chroma.Keyword:
        return theme.Blue, bg
    case chroma.Name:
        if tokenType == chroma.NameFunction {
            return theme.Yellow, bg
        }
        return theme.White, bg
    case chroma.String:
        return theme.PinkBold, bg
    case chroma.Number:
        return theme.Yellow, bg
    case chroma.Comment:
        return theme.LightGray, bg
    case chroma.Operator, chroma.Punctuation:
        return theme.Red, bg
    default:
        return theme.White, bg
    }
}