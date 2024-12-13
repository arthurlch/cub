package editor

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
)

var (
	markdownImageRegex = regexp.MustCompile(`!\[.*?\]\(.*?\)`)
	htmlImageRegex     = regexp.MustCompile(`<img.*?>`)
	plainTextFileTypes = []string{"md", "sum", "makefile", "log"}
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

	st.TextBuffer = [][]rune{}
	var lineBuffer bytes.Buffer

	fileType := strings.ToLower(filepath.Ext(filename))
	if fileType != "" {
		fileType = fileType[1:]
	}

	lexer := lexers.Match(filename)
	if lexer != nil {
		lexer = chroma.Coalesce(lexer)
		st.Language = lexer.Config().Name
	} else {
		st.Language = "Plain Text"
	}

	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		lineBuffer.Write(line)
		if !isPrefix {
			processedLine := lineBuffer.String()
			if !isPlainTextFileType(fileType) {
				processedLine = replaceImageTags(processedLine)
			}
			st.TextBuffer = append(st.TextBuffer, []rune(processedLine)) 
			lineBuffer.Reset()
		}
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

func replaceImageTags(line string) string {
	line = markdownImageRegex.ReplaceAllString(line, "[Image Placeholder]")
	line = htmlImageRegex.ReplaceAllString(line, "[Image Placeholder]")
	return line
}

func isPlainTextFileType(fileType string) bool {
	for _, plainType := range plainTextFileTypes {
		if fileType == plainType {
			return true
		}
	}
	return false
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

	for i, line := range st.TextBuffer {
		if _, err := writer.WriteString(string(line)); err != nil {
			return err
		}
		if i < len(st.TextBuffer)-1 {
			if err := writer.WriteByte('\n'); err != nil {
				return err
			}
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	st.Modified = false
	return nil
}
