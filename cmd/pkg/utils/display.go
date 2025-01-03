package utils

import (
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/syntax"
	"github.com/arthurlch/cub/cmd/pkg/theme"
	"github.com/nsf/termbox-go"
)

const LineNumberWidth = 4

// wanna void these types to avoid issues with lexer
var plainTextExtensions = map[string]bool{
    "sum": true,
    "log": true,
}

func DisplayTextBuffer(s *state.State, fileType string) {
    width, height := termbox.Size()
    if width <= LineNumberWidth {
        return 
    }
    contentWidth := width - LineNumberWidth

    if isPlainTextFile(fileType) {
        displayPlainTextBuffer(s, contentWidth, height)
        termbox.Flush()
        return
    }

    lexer := syntax.GetLexer(fileType)
    for row := 0; row < height; row++ {
        lineIndex := row + s.OffsetRow
        if lineIndex >= len(s.TextBuffer) {
            break
        }

        line := string(s.TextBuffer[lineIndex])
        renderHighlightedLine(line, row, contentWidth, lexer, s, lineIndex)
    }

    termbox.Flush()
}

func isPlainTextFile(fileType string) bool {
    return plainTextExtensions[strings.ToLower(fileType)]
}

func displayPlainTextBuffer(s *state.State, width, height int) {
    for row := 0; row < height; row++ {
        lineIndex := row + s.OffsetRow
        if lineIndex >= len(s.TextBuffer) {
            break
        }

        line := string(s.TextBuffer[lineIndex])
        startCol, endCol := getSelectionBounds(s, lineIndex, len(line))

        // Ensure we don't exceed available width
        for col := 0; col < width && col < len(line); col++ {
            var ch rune = ' '
            if col < len(line) {
                ch = rune(line[col])
            }

            fg := theme.TextForeground
            bg := theme.ColorBackground

            if startCol != -1 && col >= startCol && col < endCol {
                fg = theme.SoftBlack
                bg = theme.SelectedBackground
            }

            // Add LineNumberWidth offset to the column position
            termbox.SetCell(col + LineNumberWidth, row, ch, fg, bg)
        }

        // Fill remaining width with spaces
        for col := len(line); col < width; col++ {
            termbox.SetCell(col + LineNumberWidth, row, ' ', theme.TextForeground, theme.ColorBackground)
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
                termbox.SetCell(col + LineNumberWidth, row, rune(ch), theme.SoftBlack, theme.SelectedBackground)
            } else {
                termbox.SetCell(col + LineNumberWidth, row, rune(ch), fg, bg)
            }
            col++
        }
    }

    // Fill remaining width with spaces
    for ; col < width; col++ {
        termbox.SetCell(col + LineNumberWidth, row, ' ', theme.TextForeground, theme.ColorBackground)
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