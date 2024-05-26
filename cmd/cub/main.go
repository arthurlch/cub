package main

import (
	"fmt"
	"os"

	"github.com/arthurlch/cub/cmd/pkg/editor"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/ui"
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

func runTextEditor() {
	err := termbox.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sharedState := &state.State{} 

	editorState := &editor.EditorState{State: sharedState}
	uiState := &ui.EditorState{State: sharedState}

	if len(os.Args) > 1 {
		editorState.ReadFile(os.Args[1])
	} else {
		sharedState.TextBuffer = append(sharedState.TextBuffer, []rune{})
	}

	for {
		sharedState.Cols, sharedState.Rows = termbox.Size()
		sharedState.Rows--
		if sharedState.Cols < 78 {
			sharedState.Cols = 78
		}
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		utils.ScrollTextBuffer(sharedState)
		utils.DisplayTextBuffer(sharedState)
		uiState.StatusBar()
		termbox.SetCursor(sharedState.CurrentCol-sharedState.OffsetCol, sharedState.CurrentRow-sharedState.OffsetRow)
		termbox.Flush()
		editorState.ProcessKeyPress()
	}
}

func main() {
	runTextEditor()
}
