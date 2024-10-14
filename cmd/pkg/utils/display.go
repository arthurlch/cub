package utils

import (
	"github.com/alecthomas/chroma"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/syntax"
	"github.com/arthurlch/cub/cmd/pkg/ui"
	"github.com/nsf/termbox-go"
)

func DisplayTextBuffer(s *state.State, fileType string) {
	width, height := termbox.Size()

	for row := 0; row < height; row++ {
		if row+s.OffsetRow >= len(s.TextBuffer) {
			break // dont display if no lines
		}

		line := string(s.TextBuffer[row+s.OffsetRow])
		if syntax.IsSupportedLanguage(fileType) {
			renderHighlightedLine(line, row, width, fileType)
		} else {
			renderPlainLine(line, row, width)
		}
	}

	termbox.Flush()
}

func renderPlainLine(line string, row, width int) {
	for col, ch := range line {
		if col < width {
			termbox.SetCell(col, row, ch, ui.TextForeground, ui.ColorBackground)
		}
	}

	for col := len(line); col < width; col++ {
		termbox.SetCell(col, row, ' ', ui.TextForeground, ui.ColorBackground)
	}
}

func renderHighlightedLine(line string, row, width int, fileType string) {
	lexer := syntax.GetLexer(fileType)
	iterator, _ := lexer.Tokenise(nil, line)

	col := 0
	for token := iterator(); token != chroma.EOF; token = iterator() {
		fg, bg := syntax.GetTermboxColor(token.Type, token.Value)
		for _, ch := range token.Value {
			if col < width {
				termbox.SetCell(col, row, ch, fg, bg)
				col++
			}
		}
	}

	for ; col < width; col++ {
		termbox.SetCell(col, row, ' ', ui.TextForeground, ui.ColorBackground)
	}
}
