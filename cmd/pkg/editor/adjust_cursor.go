package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
)

func adjustCursorColToLineEnd(st *state.State) {
	if st.CurrentCol > len(st.TextBuffer[st.CurrentRow]) {
		st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
	}
	utils.Logger.Printf("adjustCursorColToLineEnd - CurrentRow: %d, CurrentCol: %d\n", st.CurrentRow, st.CurrentCol)
}
