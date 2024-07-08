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
	err := termbox.Init()
	if err != nil {
		ui.ShowErrorMessage(&state.State{}, fmt.Sprintf("Failed to initialize termbox: %v", err))
		os.Exit(1)
	}
	defer termbox.Close()

	sharedState := &state.State{}
	editorState := editor.NewEditorState(sharedState)
	uiState := ui.NewEditorState(sharedState)

	if len(os.Args) > 1 {
		err := editorState.ReadFile(os.Args[1])
		if err != nil {
			ui.ShowErrorMessage(sharedState, fmt.Sprintf("Failed to read file: %v", err))
			termbox.Flush()
			return
		}
	} else {
		sharedState.TextBuffer = append(sharedState.TextBuffer, []rune{})
	}

	mainLoop(sharedState, uiState, editorState)
}

func mainLoop(sharedState *state.State, uiState *ui.EditorState, editorState *editor.EditorState) {
	var prevCols, prevRows int

	for {
		cols, rows := termbox.Size()
		if cols != prevCols || rows != prevRows {
			sharedState.Cols, sharedState.Rows = cols, rows-1
			if sharedState.Cols < 78 {
				sharedState.Cols = 78
			}
			prevCols, prevRows = cols, rows
			redraw(sharedState, uiState)
		}

		editorState.ProcessKeyPress()
		redraw(sharedState, uiState)
	}
}

func redraw(sharedState *state.State, uiState *ui.EditorState) {
	termbox.Clear(ui.ColorBackground, termbox.ColorDefault) 
	utils.ScrollTextBuffer(sharedState)
	utils.DisplayTextBuffer(sharedState)
	uiState.StatusBar()
	termbox.SetCursor(sharedState.CurrentCol-sharedState.OffsetCol, sharedState.CurrentRow-sharedState.OffsetRow)
	termbox.Flush()
}

func main() {
	runTextEditor()
}
