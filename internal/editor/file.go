package editor

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/arthurlch/cub/internal/document"
)

var (
	mdImageRegex   = regexp.MustCompile(`!\[.*?\]\(.*?\)`)
	htmlImageRegex = regexp.MustCompile(`<img.*?>`)
)

func (e *Editor) Open(path string) error {
	e.file = path
	e.refreshLanguage()

	data, err := os.ReadFile(path)
	if err != nil {
		e.buf.SetText("")
		e.resetForFile()
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	text := string(data)
	if document.HasImageMarkup(document.Extension(path)) {
		text = replaceImageTags(text)
	}
	e.buf.SetText(text)
	e.resetForFile()
	return nil
}

func (e *Editor) Save() error {
	if e.file == "" {
		e.file = "untitled.txt"
	}
	if dir := filepath.Dir(e.file); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	if err := os.WriteFile(e.file, []byte(e.buf.String()), 0o644); err != nil {
		return err
	}
	e.modified = false
	return nil
}

func (e *Editor) resetForFile() {
	e.row, e.col = 0, 0
	e.offRow, e.offCol = 0, 0
	e.endSelection()
	e.carets = nil
	e.modified = false
	e.hist = history{}
}

func replaceImageTags(text string) string {
	text = mdImageRegex.ReplaceAllString(text, "[image]")
	return htmlImageRegex.ReplaceAllString(text, "[image]")
}
