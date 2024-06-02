package editor

import (
	"time"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/ui"
	"github.com/nsf/termbox-go"
)

func (es *EditorState) blinkCursor() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-es.State.StopBlink:
			termbox.SetCursor(es.State.CurrentCol-es.State.OffsetCol, es.State.CurrentRow-es.State.OffsetRow)
			termbox.Flush()
			return
		case <-ticker.C:
			if es.State.Mode == state.EditMode {
				es.State.Blink = !es.State.Blink
				if es.State.Blink {
					termbox.HideCursor()
				} else {
					termbox.SetCursor(es.State.CurrentCol-es.State.OffsetCol, es.State.CurrentRow-es.State.OffsetRow)
				}
			} else {
				es.State.Blink = false
				cell := termbox.CellBuffer()[(es.State.CurrentRow-es.State.OffsetRow)*es.State.Cols+(es.State.CurrentCol-es.State.OffsetCol)]
				cell.Bg = ui.CursorBackground
				termbox.SetCell(es.State.CurrentCol-es.State.OffsetCol, es.State.CurrentRow-es.State.OffsetRow, cell.Ch, cell.Fg, cell.Bg)
				termbox.SetCursor(es.State.CurrentCol-es.State.OffsetCol, es.State.CurrentRow-es.State.OffsetRow)
			}
			termbox.Flush()
		}
	}
}
