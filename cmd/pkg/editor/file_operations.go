package editor

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/arthurlch/cub/cmd/pkg/state"
)

func (es *EditorState) ReadFile(filename string) error {
	st := es.State
	if st == nil {
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
