package syntax

import "unicode"

type TokenType int

const (
    TokenKeyword TokenType = iota
    TokenString
    TokenComment
    TokenIdentifier
    TokenNumber
    TokenOperator
    TokenWhitespace
)

type Token struct {
    Type  TokenType
    Value string
}

var keywords = map[string]struct{}{
	"break": {}, "default": {}, "func": {}, "interface": {}, "select": {},
	"case": {}, "defer": {}, "go": {}, "map": {}, "struct": {},
	"chan": {}, "else": {}, "goto": {}, "package": {}, "switch": {},
	"const": {}, "fallthrough": {}, "if": {}, "range": {}, "type": {},
	"continue": {}, "for": {}, "import": {}, "return": {}, "var": {},
}

	

func Tokenize(line string) []Token {
    tokens := []Token{}
    current := ""

    for _, r := range line {
        if unicode.IsSpace(r) {
            if current != "" {
                tokens = append(tokens, createToken(current))
                current = ""
            }
            tokens = append(tokens, Token{TokenWhitespace, string(r)})
        } else if r == '"' {
            if current != "" {
                tokens = append(tokens, createToken(current))
                current = ""
            }
            current += string(r)
            tokens = append(tokens, Token{TokenString, current})
            current = ""
        } else {
            current += string(r)
        }
    }
    if current != "" {
        tokens = append(tokens, createToken(current))
    }

    return tokens
}

func createToken(current string) Token {
    if _, ok := keywords[current]; ok {
        return Token{TokenKeyword, current}
    }
    return Token{TokenIdentifier, current}
}
