package render

import (
	"fmt"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/theme"
	"github.com/nsf/termbox-go"
)

func RenderLineNumbers(st *state.State) {
	_, height := termbox.Size()

	for row := 0; row < height; row++ {
		for col := 0; col < LineNumberWidth; col++ {
			termbox.SetCell(col, row, ' ', theme.ColorDarkPink, theme.ColorBackground)
		}
	}

	for row := 0; row < height; row++ {
		lineIndex := row + st.OffsetRow
		if lineIndex >= len(st.TextBuffer) {
			break
		}

		lineNum := fmt.Sprintf("%3d ", lineIndex+1)

		for col, ch := range lineNum {
			if col < LineNumberWidth {
				termbox.SetCell(col, row, ch, theme.ColorDarkPink, theme.ColorBackground)
			}
		}
	}
}
