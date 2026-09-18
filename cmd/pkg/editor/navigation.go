package editor

/*
Notes:
The direction concerning the navigation is tough.
I decided to go toward more a kakoune like navigation and adding a little of vim motion as well
So far I think it goes well without being too complicated.

There is some design changes we could think about like for line jump we could have some prompt.

I don't wish to add much more feature when it comes to the navigation.
The philophy is to have better navigation than nano but still keeping it simple.

As for the mode handling we will keep the separation between simple nav and complex nav.

It's easy to handle and match well with the only 2 modes.
*/

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/arthurlch/cub/cmd/pkg/logging"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/nsf/termbox-go"
)

func handleNavigation(st *state.State, keyEvent termbox.Event) {
	handleSimpleNavigation(st, keyEvent)

	st.ClampCursor()
	logging.LogKeyPress("handleNavigation", keyEvent)
	logging.LogBufferState(st, "Navigation")
}

func handleSimpleNavigation(st *state.State, keyEvent termbox.Event) {
	switch keyEvent.Key {
	case termbox.KeyArrowUp:
		moveUp(st)
	case termbox.KeyArrowDown:
		moveDown(st)
	case termbox.KeyArrowLeft:
		moveLeft(st)
	case termbox.KeyArrowRight:
		moveRight(st)
	case termbox.KeyHome:
		st.CurrentCol = 0
	case termbox.KeyEnd:
		st.SnapCursorToLineEnd()
	case termbox.KeyPgup:
		movePageUp(st)
	case termbox.KeyPgdn:
		movePageDown(st)
	}
	st.ClampCursor()
}

func moveUp(st *state.State) {
	st.MoveCursor(st.CurrentRow-1, st.CurrentCol)
}

func moveDown(st *state.State) {
	st.MoveCursor(st.CurrentRow+1, st.CurrentCol)
}

func moveLeft(st *state.State) {
	if st.CurrentCol > 0 {
		st.MoveCursor(st.CurrentRow, st.CurrentCol-1)
	} else if st.CurrentRow > 0 {
		st.MoveCursor(st.CurrentRow-1, len(st.TextBuffer[st.CurrentRow-1]))
	}
}

func moveRight(st *state.State) {
	if st.CurrentCol < len(st.TextBuffer[st.CurrentRow]) {
		st.MoveCursor(st.CurrentRow, st.CurrentCol+1)
	} else if st.CurrentRow < len(st.TextBuffer)-1 {
		st.MoveCursor(st.CurrentRow+1, 0)
	}
}

func movePageUp(st *state.State) {
	st.MoveCursor(st.CurrentRow-int(st.Rows/4), st.CurrentCol)
}

func movePageDown(st *state.State) {
	st.MoveCursor(st.CurrentRow+int(st.Rows/4), st.CurrentCol)
}

func moveToNextWord(st *state.State) {
	for row := st.CurrentRow; row < len(st.TextBuffer); row++ {
		line := st.TextBuffer[row]
		for col := st.CurrentCol + 1; col < len(line); col++ {
			if isWordBoundary(line[col-1], line[col]) {
				st.MoveCursor(row, col)
				return
			}
		}
		st.CurrentCol = 0
	}
}

func moveToPreviousWord(st *state.State) {
	for row := st.CurrentRow; row >= 0; row-- {
		line := st.TextBuffer[row]
		for col := st.CurrentCol - 1; col > 0; col-- {
			if isWordBoundary(line[col-1], line[col]) {
				st.MoveCursor(row, col)
				return
			}
		}
		st.CurrentCol = len(line)
	}
}

func isWordBoundary(prev, next rune) bool {
	return unicode.IsSpace(prev) && !unicode.IsSpace(next)
}

func moveToMatchingBracket(st *state.State, openBracket rune) {
	matching := map[rune]rune{'(': ')', ')': '(', '{': '}', '}': '{'}
	closeBracket := matching[openBracket]
	depth := 0

	if openBracket == '(' || openBracket == '{' {
		for row := st.CurrentRow; row < len(st.TextBuffer); row++ {
			line := st.TextBuffer[row]
			for col := 0; col < len(line); col++ {
				switch line[col] {
				case openBracket:
					depth++
				case closeBracket:
					depth--
					if depth == 0 {
						st.MoveCursor(row, col)
						return
					}
				}
			}
		}
	} else {
		for row := st.CurrentRow; row >= 0; row-- {
			line := st.TextBuffer[row]
			for col := len(line) - 1; col >= 0; col-- {
				switch line[col] {
				case closeBracket:
					depth++
				case openBracket:
					depth--
					if depth == 0 {
						st.MoveCursor(row, col)
						return
					}
				}
			}
		}
	}
}

func moveToNextEmptyLine(st *state.State) {
	for row := st.CurrentRow + 1; row < len(st.TextBuffer); row++ {
		if strings.TrimSpace(string(st.TextBuffer[row])) == "" {
			st.MoveCursor(row, 0)
			return
		}
	}
}

func moveToPreviousEmptyLine(st *state.State) {
	for row := st.CurrentRow - 1; row >= 0; row-- {
		if strings.TrimSpace(string(st.TextBuffer[row])) == "" {
			st.MoveCursor(row, 0)
			return
		}
	}
}

func moveToLineStart(st *state.State) {
	line := st.TextBuffer[st.CurrentRow]
	for col := 0; col < len(line); col++ {
		if !unicode.IsSpace(line[col]) {
			st.MoveCursor(st.CurrentRow, col)
			return
		}
	}
	st.MoveCursor(st.CurrentRow, 0)
}

func moveToLineEnd(st *state.State) {
	st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
}

func jumpToLine(st *state.State) {
	lineNumber, err := strconv.Atoi(st.LineNumberBuffer)
	if err != nil {
		return
	}
	// end of the buffer is max, I shall not jump further
	if lineNumber > len(st.TextBuffer) {
		lineNumber = len(st.TextBuffer)
	}

	st.MoveCursor(lineNumber-1, 0)
}
