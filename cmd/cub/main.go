package main

import (
	"fmt"
	"os"

	"github.com/arthurlch/cub/cmd/pkg/document"
	"github.com/arthurlch/cub/cmd/pkg/editor"
	"github.com/arthurlch/cub/cmd/pkg/logging"
	"github.com/arthurlch/cub/cmd/pkg/render"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/theme"
	"github.com/nsf/termbox-go"
)

var version = "dev"

const usage = `cub - a simple terminal text editor

usage:
  cub [file]

flags:
  -h, --help      show this help and exit
  -v, --version   show version and exit

keys:
  i              enter insert mode        Esc            back to view mode
  Ctrl+S         save                     Ctrl+Q         quit
  Ctrl+U         undo                     Ctrl+R         redo
  Ctrl+H         help
`

func runTextEditor() {
	logging.InitLogger()

	if err := termbox.Init(); err != nil {
		fmt.Println("Failed to initialize termbox:", err)
		os.Exit(1)
	}
	defer termbox.Close()

	if termbox.SetOutputMode(termbox.OutputRGB) != termbox.OutputRGB {
		fmt.Println("Failed to enable RGB mode.")
		os.Exit(1)
	}

	sharedState := &state.State{}
	editorState := editor.NewEditorState(sharedState)

	var fileType string
	if len(os.Args) > 1 {
		filePath := os.Args[1]
		if err := editorState.ReadFile(filePath); err != nil {
			sharedState.ShowMessage(fmt.Sprintf("Failed to read file: %v", err))
			return
		}
		fileType = document.Extension(filePath)
	} else {
		sharedState.TextBuffer = append(sharedState.TextBuffer, []rune{})
	}

	mainLoop(sharedState, editorState, fileType)
}

func mainLoop(sharedState *state.State, editorState *editor.EditorState, fileType string) {
	sharedState.Cols, sharedState.Rows = termbox.Size()
	sharedState.Rows--

	for !sharedState.Quit {
		redraw(sharedState, fileType)
		editorState.ProcessKeyPress(fileType)
	}
}

func redraw(sharedState *state.State, fileType string) {
	_ = termbox.Clear(theme.TextForeground, theme.ColorBackground)

	sharedState.Scroll()

	render.DisplayTextBuffer(sharedState, fileType)

	render.RenderLineNumbers(sharedState)

	render.StatusBar(sharedState)

	cursorCol := sharedState.CurrentCol - sharedState.OffsetCol + render.LineNumberWidth
	cursorRow := sharedState.CurrentRow - sharedState.OffsetRow
	termbox.SetCursor(cursorCol, cursorRow)

	_ = termbox.Flush()
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version":
			fmt.Println("cub", version)
			return
		case "-h", "--help":
			fmt.Print(usage)
			return
		}
	}
	runTextEditor()
}
