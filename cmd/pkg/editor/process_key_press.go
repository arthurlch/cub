package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

func (es *EditorState) ProcessKeyPress(fileType string) {
	st := es.State
	keyEvent := utils.GetKey()

	switch keyEvent.Type {
	case termbox.EventKey:
		handleKeyPress(es, keyEvent)
	case termbox.EventResize:
		st.Cols, st.Rows = termbox.Size()
		st.Rows--
	}
}
