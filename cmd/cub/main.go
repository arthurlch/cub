package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/ui"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type EditorState struct {
	state *state.State
}

func NewEditorState() *EditorState {
	return &EditorState{state: &state.State{}}
}

func (es *EditorState) readFile(filename string) {
	st := es.state
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

func displayTextBuffer(s *state.State) {
	var row, col int
	for row = 0; row < s.Rows; row++ {
		textBufferRow := row + s.OffsetRow
		for col = 0; col < s.Cols; col++ {
			textBufferCol := col + s.OffsetCol
			if textBufferRow >= 0 && textBufferRow < len(s.TextBuffer) && textBufferCol < len(s.TextBuffer[textBufferRow]) {
				if s.TextBuffer[textBufferRow][textBufferCol] != '\t' {
					termbox.SetChar(col, row, s.TextBuffer[textBufferRow][textBufferCol])
				} else {
					termbox.SetCell(col, row, rune(' '), termbox.ColorDefault, termbox.ColorDefault)
				}
			} else if row+s.OffsetCol > len(s.TextBuffer) {
				termbox.SetCell(0, row, rune('*'), termbox.ColorLightMagenta, termbox.ColorDefault)
				termbox.SetChar(col, row, rune('\n'))
			}
		}
	}
}

func scrollTextBuffer(s *state.State) {
	if s.CurrentRow < s.OffsetRow {
		s.OffsetRow = s.CurrentRow
	}

	if s.CurrentRow >= s.OffsetRow+s.Rows {
		s.OffsetRow = s.CurrentRow - s.Rows + 1
	}

	if s.CurrentCol < s.OffsetCol {
		s.OffsetCol = s.CurrentCol
	}

	if s.CurrentCol >= s.OffsetCol+s.Cols {
		s.OffsetCol = s.CurrentCol - s.Cols + 1
		if s.OffsetCol < 0 {
			s.OffsetCol = 0
		}
	}
}

func getKey() termbox.Event {
	var key_event termbox.Event
	switch event := termbox.PollEvent(); event.Type {
	case termbox.EventKey:
		key_event = event
	case termbox.EventError:
		panic(event.Err)
	}
	return key_event
}

func (es *EditorState) processKeyPress() {
	st := es.state
	keyEvent := getKey()
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
		scrollTextBuffer(st)
	} else if st.Mode == state.EditMode {
		if keyEvent.Ch != 0 {
			if st.CurrentRow >= len(st.TextBuffer) {
				return
			}
			es.insertRunes(keyEvent)
			st.Modified = true
		} else {
			switch keyEvent.Key {
			case termbox.KeyEnter:
				es.insertNewLine()
				st.Modified = true
			case termbox.KeyArrowUp, termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight, termbox.KeyHome, termbox.KeyEnd, termbox.KeyPgup, termbox.KeyPgdn:
				break
			case termbox.KeyTab, termbox.KeySpace:
				for i := 0; i < 4; i++ {
					es.insertRunes(keyEvent)
				}
				st.Modified = true
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				es.deleteRune()
				st.Modified = true
			}

			if st.CurrentCol > len(st.TextBuffer[st.CurrentRow]) {
				st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
			}
		}
		scrollTextBuffer(st)
	}
}

func (es *EditorState) insertRunes(event termbox.Event) {
	st := es.state
	if st.CurrentRow >= len(st.TextBuffer) {
		return
	}

	newRow := make([]rune, len(st.TextBuffer[st.CurrentRow])+1)
	copy(newRow, st.TextBuffer[st.CurrentRow][:st.CurrentCol])

	switch event.Key {
	case termbox.KeySpace, termbox.KeyTab:
		newRow[st.CurrentCol] = ' '
	default:
		newRow[st.CurrentCol] = event.Ch
	}

	copy(newRow[st.CurrentCol+1:], st.TextBuffer[st.CurrentRow][st.CurrentCol:])
	st.TextBuffer[st.CurrentRow] = newRow
	st.CurrentCol++
}

func (es *EditorState) deleteRune() {
	st := es.state
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

func (es *EditorState) insertNewLine() {
	st := es.state
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

func print_message(col, row int, foreground, background termbox.Attribute, message string) {
	for _, ch := range message {
		termbox.SetCell(col, row, ch, foreground, background)
		col += runewidth.RuneWidth(ch)
	}
}

func runTextEditor() {
	err := termbox.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	editorState := NewEditorState()
	uiState := ui.NewEditorState() 
	uiState.State = editorState.state

	if len(os.Args) > 1 {
		editorState.readFile(os.Args[1])
		uiEditorState.State = editorState.state 
	} else {
		editorState.state.TextBuffer = append(editorState.state.TextBuffer, []rune{})
		uiEditorState.State = editorState.state 
	}

	for {
		editorState.state.Cols, editorState.state.Rows = termbox.Size()
		editorState.state.Rows--
		if editorState.state.Cols < 78 {
			editorState.state.Cols = 78
		}
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		scrollTextBuffer(editorState.state)
		displayTextBuffer(editorState.state)
		uiEditorState.StatusBar() 
		termbox.SetCursor(editorState.state.CurrentCol-editorState.state.OffsetCol, editorState.state.CurrentRow-editorState.state.OffsetRow)
		termbox.Flush()
		editorState.processKeyPress()
	}
}



func main() {
	runTextEditor()
}
