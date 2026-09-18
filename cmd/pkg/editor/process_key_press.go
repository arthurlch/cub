package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/input"
	"github.com/arthurlch/cub/cmd/pkg/logging"
	"github.com/nsf/termbox-go"
)

func (es *EditorState) ProcessKeyPress(fileType string) {
	st := es.State
	event := input.PollEvent()

	switch event.Type {
	case termbox.EventKey:
		handleKeyPress(es, event)
	case termbox.EventResize:
		st.Cols, st.Rows = termbox.Size()
		st.Rows--
	case termbox.EventError:
		logging.Logger.Printf("input error: %v", event.Err)
	}
}
