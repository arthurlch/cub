package document

import (
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

var noHighlight = map[string]bool{
	"sum": true,
	"log": true,
}

var imageMarkup = map[string]bool{
	"md":       true,
	"markdown": true,
	"html":     true,
	"htm":      true,
	"rst":      true,
}

func Extension(path string) string {
	return strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
}

func IsPlainText(fileType string) bool {
	return noHighlight[strings.ToLower(fileType)]
}

func HasImageMarkup(fileType string) bool {
	return imageMarkup[strings.ToLower(fileType)]
}

func Language(path string) string {
	lexer := lexers.Match(path)
	if lexer == nil {
		return "Plain Text"
	}
	return chroma.Coalesce(lexer).Config().Name
}
