package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arthurlch/cub/cmd/pkg/editor"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/ui"
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

var logger *log.Logger

func init() {
	file, err := os.OpenFile("editor.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}
	logger = log.New(file, "", log.LstdFlags)
}

func runTextEditor() {
	utils.InitLogger()

	if err := termbox.Init(); err != nil {
		ui.ShowErrorMessage(&state.State{}, fmt.Sprintf("Failed to initialize termbox: %v", err))
		os.Exit(1)
	}
	defer termbox.Close()

	if termbox.SetOutputMode(termbox.OutputRGB) != termbox.OutputRGB {
		fmt.Println("Failed to enable RGB mode.")
		os.Exit(1)
	}

	sharedState := &state.State{}
	editorState := editor.NewEditorState(sharedState)
	uiState := ui.NewEditorState(sharedState)

	var fileType string
	if len(os.Args) > 1 {
		filePath := os.Args[1]
		if err := editorState.ReadFile(filePath); err != nil {
			ui.ShowErrorMessage(sharedState, fmt.Sprintf("Failed to read file: %v", err))
			termbox.Flush()
			return
		}
		fileType = utils.DetermineFileType(filePath)
	} else {
		sharedState.TextBuffer = append(sharedState.TextBuffer, []rune{})
		fileType = ""
	}

	mainLoop(sharedState, uiState, editorState, fileType)
}

func mainLoop(sharedState *state.State, uiState *ui.EditorState, editorState *editor.EditorState, fileType string) {
	var prevCols, prevRows int

	for {
		cols, rows := termbox.Size()
		if cols != prevCols || rows != prevRows {
			sharedState.Cols, sharedState.Rows = cols, rows-1
			prevCols, prevRows = cols, rows
			redraw(sharedState, uiState, fileType)
		}

		editorState.ProcessKeyPress(fileType)
		redraw(sharedState, uiState, fileType)
	}
}

// filetype not used but likely to be used later on
func redraw(sharedState *state.State, uiState *ui.EditorState, fileType string) {
	termbox.Clear(ui.TextForeground, ui.ColorBackground)
	utils.ScrollTextBuffer(sharedState)
	utils.DisplayTextBuffer(sharedState, fileType) 
	uiState.StatusBar()
	termbox.SetCursor(sharedState.CurrentCol-sharedState.OffsetCol, sharedState.CurrentRow-sharedState.OffsetRow)
	termbox.Flush() 
	uiState.StatusBar()
	termbox.SetCursor(sharedState.CurrentCol-sharedState.OffsetCol, sharedState.CurrentRow-sharedState.OffsetRow)
	termbox.Flush()
}


func main() {
	runTextEditor()
}
