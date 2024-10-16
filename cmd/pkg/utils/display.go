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
			break // don't display if no lines
		}

		lineIndex := row + s.OffsetRow
		line := string(s.TextBuffer[lineIndex])

		if syntax.IsSupportedLanguage(fileType) {
			renderHighlightedLine(line, row, width, fileType, s, lineIndex)
		} else {
			renderPlainLine(line, row, width, s, lineIndex)
		}
	}

	termbox.Flush()
}

func renderPlainLine(line string, row, width int, s *state.State, lineIndex int) {
	startCol, endCol := -1, -1

	if s.SelectionActive {
		if lineIndex > s.StartRow && lineIndex < s.EndRow {
			startCol, endCol = 0, len(line)
		} else if lineIndex == s.StartRow && lineIndex == s.EndRow {
			startCol, endCol = s.StartCol, s.EndCol
		} else if lineIndex == s.StartRow {
			startCol = s.StartCol
			endCol = len(line)
		} else if lineIndex == s.EndRow {
			startCol = 0
			endCol = s.EndCol
		}
	}

	for col, ch := range line {
		if startCol != -1 && col >= startCol && col < endCol {
			termbox.SetCell(col, row, ch, ui.SoftBlack, ui.SelectedBackground)
		} else {
			termbox.SetCell(col, row, ch, ui.TextForeground, ui.ColorBackground)
		}
	}

	for col := len(line); col < width; col++ {
		termbox.SetCell(col, row, ' ', ui.TextForeground, ui.ColorBackground)
	}
}

func renderHighlightedLine(line string, row, width int, fileType string, s *state.State, lineIndex int) {
	lexer := syntax.GetLexer(fileType)
	iterator, _ := lexer.Tokenise(nil, line)

	startCol, endCol := -1, -1
	if s.SelectionActive {
		if lineIndex > s.StartRow && lineIndex < s.EndRow {
			startCol, endCol = 0, len(line)
		} else if lineIndex == s.StartRow && lineIndex == s.EndRow {
			startCol, endCol = s.StartCol, s.EndCol
		} else if lineIndex == s.StartRow {
			startCol = s.StartCol
			endCol = len(line)
		} else if lineIndex == s.EndRow {
			startCol = 0
			endCol = s.EndCol
		}
	}

	col := 0
	for token := iterator(); token != chroma.EOF; token = iterator() {
		fg, bg := syntax.GetTermboxColor(token.Type, token.Value)
		for _, ch := range token.Value {
			if col >= startCol && col < endCol {
				termbox.SetCell(col, row, ch, ui.SoftBlack, ui.White)
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
