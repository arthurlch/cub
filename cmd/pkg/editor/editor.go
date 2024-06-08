package editor

import (
	"bufio"
	"log"
	"os"
	"path/filepath"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/nsf/termbox-go"
)

type EditorState struct {
	State *state.State
}

func NewEditorState(sharedState *state.State) *EditorState {
	if sharedState == nil {
		log.Println("Shared state is nil, creating a new state.")
		sharedState = &state.State{}
	}
	return &EditorState{State: sharedState}
}

func (es *EditorState) ReadFile(filename string) error {
	st := es.State
	if st == nil {
		log.Println("Editor state is nil, initializing.")
		st = &state.State{}
		es.State = st
	}

	file, err := os.Open(filename)
	if err != nil {
		st.SourceFile = filename
		st.TextBuffer = append(st.TextBuffer, []rune{})
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		st.TextBuffer = append(st.TextBuffer, []rune(line))
	}
	if len(st.TextBuffer) == 0 {
		st.TextBuffer = append(st.TextBuffer, []rune{})
	}
	st.SourceFile = filename
	return scanner.Err()
}

func (es *EditorState) SaveFile() error {
	st := es.State
	if st == nil {
		return os.ErrInvalid
	}

	filename := st.SourceFile
	if filename == "" {
		filename = "untitled.txt"
		st.SourceFile = filename
	}

	// Ensure the directory exists
	dir := filepath.Dir(filename)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range st.TextBuffer {
		_, err := writer.WriteString(string(line) + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush()
	st.Modified = false
	return nil
}

func (es *EditorState) InsertRunes(event termbox.Event) {
	st := es.State
	if st == nil || st.CurrentRow >= len(st.TextBuffer) {
		return
	}

	row := st.TextBuffer[st.CurrentRow]
	if st.CurrentCol > len(row) {
		st.CurrentCol = len(row)
	}

	newRow := make([]rune, len(row)+1)
	copy(newRow, row[:st.CurrentCol])

	if event.Key == termbox.KeySpace {
		newRow[st.CurrentCol] = ' '
	} else {
		newRow[st.CurrentCol] = event.Ch
	}

	copy(newRow[st.CurrentCol+1:], row[st.CurrentCol:])
	st.TextBuffer[st.CurrentRow] = newRow
	st.CurrentCol++
}

func (es *EditorState) DeleteRune() {
	st := es.State
	if st == nil || st.CurrentRow >= len(st.TextBuffer) {
		return
	}

	row := st.TextBuffer[st.CurrentRow]
	if st.CurrentCol > 0 {
		st.CurrentCol--
		if len(row) > 0 {
			st.TextBuffer[st.CurrentRow] = append(row[:st.CurrentCol], row[st.CurrentCol+1:]...)
		}
	} else if st.CurrentRow > 0 {
		prevLineLen := len(st.TextBuffer[st.CurrentRow-1])
		st.TextBuffer[st.CurrentRow-1] = append(st.TextBuffer[st.CurrentRow-1], row...)
		st.TextBuffer = append(st.TextBuffer[:st.CurrentRow], st.TextBuffer[st.CurrentRow+1:]...)
		st.CurrentRow--
		st.CurrentCol = prevLineLen
	}
}

func (es *EditorState) InsertNewLine() {
	st := es.State
	if st == nil || st.CurrentRow >= len(st.TextBuffer) {
		return
	}
	beforeCursor := st.TextBuffer[st.CurrentRow][:st.CurrentCol]
	afterCursor := st.TextBuffer[st.CurrentRow][st.CurrentCol:]

	st.TextBuffer[st.CurrentRow] = beforeCursor

	if st.CurrentRow+1 < len(st.TextBuffer) {
		st.TextBuffer = append(st.TextBuffer[:st.CurrentRow+1], append([][]rune{afterCursor}, st.TextBuffer[st.CurrentRow+1:]...)...)
	} else {
		st.TextBuffer = append(st.TextBuffer, afterCursor)
	}

	st.CurrentRow++
	st.CurrentCol = 0
}