package editor

import (
	"bufio"
	"os"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

type EditorState struct {
	State *state.State
}

func NewEditorState() *EditorState {
	return &EditorState{State: &state.State{}}
}

func (es *EditorState) ReadFile(filename string) {
	st := es.State
	file, err := os.Open(filename)
	if err != nil {
		st.SourceFile = filename
		st.TextBuffer = append(st.TextBuffer, []rune{})
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		line := scanner.Text()
		st.TextBuffer = append(st.TextBuffer, []rune{})
		for i := 0; i < len(line); i++ {
			st.TextBuffer[lineNumber] = append(st.TextBuffer[lineNumber], rune(line[i]))
		}
		lineNumber++
	}
	if lineNumber == 0 {
		st.TextBuffer = append(st.TextBuffer, []rune{})
	}
}

func (es *EditorState) InsertRunes(event termbox.Event) {
	st := es.State
	if st.CurrentRow >= len(st.TextBuffer) {
		return
	}

	newRow := make([]rune, len(st.TextBuffer[st.CurrentRow])+1)
	copy(newRow, st.TextBuffer[st.CurrentRow][:st.CurrentCol])

	if event.Key == termbox.KeySpace {
		newRow[st.CurrentCol] = ' '
	} else {
		newRow[st.CurrentCol] = event.Ch
	}

	copy(newRow[st.CurrentCol+1:], st.TextBuffer[st.CurrentRow][st.CurrentCol:])
	st.TextBuffer[st.CurrentRow] = newRow
	st.CurrentCol++
}

func (es *EditorState) DeleteRune() {
	st := es.State
	if st.CurrentCol > 0 && st.CurrentRow < len(st.TextBuffer) {
		st.CurrentCol--
		st.TextBuffer[st.CurrentRow] = append(st.TextBuffer[st.CurrentRow][:st.CurrentCol], st.TextBuffer[st.CurrentRow][st.CurrentCol+1:]...)
	} else if st.CurrentRow > 0 {
		prevLineLen := len(st.TextBuffer[st.CurrentRow-1])
		st.TextBuffer[st.CurrentRow-1] = append(st.TextBuffer[st.CurrentRow-1], st.TextBuffer[st.CurrentRow]...)
		st.TextBuffer = append(st.TextBuffer[:st.CurrentRow], st.TextBuffer[st.CurrentRow+1:]...)
		st.CurrentRow--
		st.CurrentCol = prevLineLen
	}
}

func (es *EditorState) InsertNewLine() {
	st := es.State
	if st.CurrentRow >= len(st.TextBuffer) {
		return
	}
	beforeCursor := st.TextBuffer[st.CurrentRow][:st.CurrentCol]
	afterCursor := make([]rune, len(st.TextBuffer[st.CurrentRow][st.CurrentCol:]))
	copy(afterCursor, st.TextBuffer[st.CurrentRow][st.CurrentCol:])

	st.TextBuffer[st.CurrentRow] = beforeCursor

	if st.CurrentRow+1 < len(st.TextBuffer) {
		st.TextBuffer = append(st.TextBuffer[:st.CurrentRow+1], append([][]rune{afterCursor}, st.TextBuffer[st.CurrentRow+1:]...)...)
	} else {
		st.TextBuffer = append(st.TextBuffer, afterCursor)
	}

	st.CurrentRow++
	st.CurrentCol = 0
}

func (es *EditorState) ProcessKeyPress() {
	st := es.State
	keyEvent := utils.GetKey()
	if keyEvent.Key == termbox.KeyEsc {
		st.Mode = state.ViewMode
		st.QuitKey = termbox.KeyEsc
		return
	}

	if st.QuitKey == termbox.KeyEsc && keyEvent.Ch == 'q' {
		termbox.Close()
		os.Exit(0)
	} else {
		st.QuitKey = 0
	}

	if keyEvent.Ch == 'e' {
		st.Mode = state.EditMode
		return
	}

	if st.Mode == state.ViewMode {
		switch keyEvent.Key {
		case termbox.KeyArrowUp:
			if st.CurrentRow != 0 {
				st.CurrentRow--
			}
		case termbox.KeyArrowDown:
			if st.CurrentRow < len(st.TextBuffer)-1 {
				st.CurrentRow++
			}
		case termbox.KeyArrowLeft:
			if st.CurrentCol != 0 {
				st.CurrentCol--
			} else if st.CurrentRow > 0 {
				st.CurrentRow--
				st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
			}
		case termbox.KeyArrowRight:
			if st.CurrentCol < len(st.TextBuffer[st.CurrentRow]) {
				st.CurrentCol++
			} else if st.CurrentRow < len(st.TextBuffer)-1 {
				st.CurrentRow++
				st.CurrentCol = 0
			}
		case termbox.KeyHome:
			st.CurrentCol = 0
		case termbox.KeyEnd:
			st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
		case termbox.KeyPgup:
			if st.CurrentRow-int(st.Rows/4) > 0 {
				st.CurrentRow -= int(st.Rows / 4)
			}
		case termbox.KeyPgdn:
			if st.CurrentRow+int(st.Rows/4) < len(st.TextBuffer)-1 {
				st.CurrentRow += int(st.Rows / 4)
			}
		}
		utils.ScrollTextBuffer(st)
	} else if st.Mode == state.EditMode {
		if keyEvent.Ch != 0 {
			if st.CurrentRow >= len(st.TextBuffer) {
				return
			}
			es.InsertRunes(keyEvent)
			st.Modified = true
		} else {
			switch keyEvent.Key {
			case termbox.KeyEnter:
				es.InsertNewLine()
				st.Modified = true
			case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight, termbox.KeyHome, termbox.KeyEnd, termbox.KeyPgup, termbox.KeyPgdn:
				break
			case termbox.KeyTab:
				for i := 0; i < 4; i++ {
					es.InsertRunes(keyEvent)
				}
				st.Modified = true
			case termbox.KeySpace:
				es.InsertRunes(keyEvent)
				st.Modified = true
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				es.DeleteRune()
				st.Modified = true
			}

			if st.CurrentCol > len(st.TextBuffer[st.CurrentRow]) {
				st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
			}
		}
		utils.ScrollTextBuffer(st)
	}
}
