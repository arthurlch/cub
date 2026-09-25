package syntax

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

func GetLexer(fileType string) chroma.Lexer {
	lexer := lexers.Get(fileType)
	if lexer == nil {
		return lexers.Fallback
	}
	return chroma.Coalesce(lexer)
}
