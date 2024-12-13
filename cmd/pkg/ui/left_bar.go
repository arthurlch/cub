package ui

import (
	"fmt"

	"github.com/arthurlch/cub/cmd/pkg/theme"
	"github.com/nsf/termbox-go"
)

const LineNumberWidth = 4

func (es *EditorState) RenderLineNumbers() {
	st := es.State
	maxLines := len(st.TextBuffer)
	numWidth := 1
	if maxLines > 9 {
			numWidth = 2
	}
	if maxLines > 99 {
			numWidth = 3
	}
	if maxLines > 999 {
			numWidth = 4
	}
	if maxLines > 9999 {
			numWidth = 5
	}
	actualWidth := numWidth + 1

	for i := 0; i < st.Rows; i++ {
			for j := 0; j < actualWidth; j++ {
					termbox.SetCell(j, i, ' ', theme.ColorDarkPink, theme.ColorBackground)
			}
			
			if i < len(st.TextBuffer) {
					lineNumber := fmt.Sprintf("%*d", numWidth, i+st.OffsetRow+1)  // Right align number
					printMessage(0, i, theme.ColorDarkPink, theme.ColorBackground, lineNumber)
			}
	}
}