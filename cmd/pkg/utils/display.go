package utils

import (
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/syntax"
	"github.com/arthurlch/cub/cmd/pkg/ui"
	"github.com/nsf/termbox-go"
)

var plainTextExtensions = map[string]bool{
	"md":       true,
	"sum":      true,
	"makefile": true,
	"lock":     true,
	"txt":      true,
	"log":      true,
	"conf":     true,
	"config":   true,
	"yml":      true,
	"yaml":     true,
	"json":     true,
	"toml":     true,
	"env":      true,
}

func DisplayTextBuffer(s *state.State, fileType string) {
	width, height := termbox.Size()

	if isPlainTextFile(fileType) {
		displayPlainTextBuffer(s, width, height)
		termbox.Flush()
		return
	}

	lexer := syntax.GetLexer(fileType)

	for row := 0; row < height; row++ {
		if row+s.OffsetRow >= len(s.TextBuffer) {
			break
		}

		lineIndex := row + s.OffsetRow
		line := string(s.TextBuffer[lineIndex])
		renderHighlightedLine(line, row, width, lexer, s, lineIndex)
	}

	termbox.Flush()
}

func isPlainTextFile(fileType string) bool {
	fileType = strings.ToLower(fileType)
	return plainTextExtensions[fileType]
}

func displayPlainTextBuffer(s *state.State, width, height int) {
	for row := 0; row < height; row++ {
		if row+s.OffsetRow >= len(s.TextBuffer) {
			break
		}

		lineIndex := row + s.OffsetRow
		line := string(s.TextBuffer[lineIndex])
		startCol, endCol := getSelectionBounds(s, lineIndex, len(line))

		for col := 0; col < width; col++ {
			var ch rune = ' '
			if col < len(line) {
				ch = rune(line[col])
			}

			fg := ui.TextForeground
			bg := ui.ColorBackground

			if startCol != -1 && col >= startCol && col < endCol {
				fg = ui.SoftBlack
				bg = ui.SelectedBackground
			}

			termbox.SetCell(col, row, ch, fg, bg)
		}
	}
}

func renderHighlightedLine(line string, row, width int, lexer chroma.Lexer, s *state.State, lineIndex int) {
	iterator, _ := lexer.Tokenise(nil, line)
	startCol, endCol := getSelectionBounds(s, lineIndex, len(line))

	col := 0
	for token := iterator(); token != chroma.EOF; token = iterator() {
		fg, bg := syntax.GetTermboxColor(token.Type, token.Value)
		for _, ch := range token.Value {
			if col >= width {
				return
			}

			if col >= startCol && col < endCol {
				termbox.SetCell(col, row, ch, ui.SoftBlack, ui.SelectedBackground)
			} else {
				termbox.SetCell(col, row, ch, fg, bg)
			}
			col++
		}
	}

	for ; col < width; col++ {
		termbox.SetCell(col, row, ' ', ui.TextForeground, ui.ColorBackground)
	}
}

func getSelectionBounds(s *state.State, lineIndex, lineLength int) (startCol, endCol int) {
	if !s.SelectionActive {
		return -1, -1
	}

	if lineIndex > s.StartRow && lineIndex < s.EndRow {
		return 0, lineLength
	} else if lineIndex == s.StartRow && lineIndex == s.EndRow {
		return s.StartCol, s.EndCol
	} else if lineIndex == s.StartRow {
		return s.StartCol, lineLength
	} else if lineIndex == s.EndRow {
		return 0, s.EndCol
	}

	return -1, -1
}