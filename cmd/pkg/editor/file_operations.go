package editor

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
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

	const bufferSize = 64 * 1024 // 64KB buffer
	reader := bufio.NewReaderSize(file, bufferSize)
	
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, reader)
	if err != nil {
		return err
	}

	lines := bytes.Split(buffer.Bytes(), []byte{'\n'})
	
	st.TextBuffer = make([][]rune, 0, len(lines))
	
	for _, line := range lines {
		st.TextBuffer = append(st.TextBuffer, bytes.Runes(line))
	}

	if len(st.TextBuffer) == 0 {
		st.TextBuffer = append(st.TextBuffer, []rune{})
	}
	st.SourceFile = filename
	
	st.UndoBuffer = append(st.UndoBuffer, state.UndoState{
		TextBuffer: utils.DeepCopyTextBuffer(st.TextBuffer),
		CurrentRow: st.CurrentRow,
		CurrentCol: st.CurrentCol,
	})

	return nil
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

	const writeBufferSize = 64 * 1024 // 64KB buffer
	writer := bufio.NewWriterSize(file, writeBufferSize)
	
	var buffer bytes.Buffer
	buffer.Grow(1024) 

	for i, line := range st.TextBuffer {
		buffer.WriteString(string(line))
		if i < len(st.TextBuffer)-1 {
			buffer.WriteByte('\n')
		}
	}

	_, err = writer.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	st.Modified = false
	return nil
}