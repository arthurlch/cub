// cmd/pkg/editor/buffer_operation_test.go

package editor

import (
	"testing"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/stretchr/testify/assert"
)

func TestUndo(t *testing.T) {
	st := &state.State{
		TextBuffer: []([]rune){
			[]rune("Hello, World!"),
		},
		HistoryIndex: 1,
		ChangeHistory: []state.Change{
			{
				Type: state.Insert,
				Row:  0,
				Col:  12,
				Text: []rune("!"),
			},
		},
		CurrentRow: 0,
		CurrentCol: 13,
	}

	Undo(st)

	assert.Equal(t, 0, st.HistoryIndex, "HistoryIndex should decrement after Undo")
	assert.Equal(t, []rune("Hello, World"), st.TextBuffer[0], "TextBuffer should have the last change undone")
	assert.Equal(t, 12, st.CurrentCol, "CurrentCol should reflect the undone change (move back)")
}

func TestRedo(t *testing.T) {
	st := &state.State{
		TextBuffer: []([]rune){
			[]rune("Hello, World"),
		},
		HistoryIndex: 0,
		ChangeHistory: []state.Change{
			{
				Type: state.Insert,
				Row:  0,
				Col:  12,
				Text: []rune("!"),
			},
		},
		CurrentRow: 0,
		CurrentCol: 12,
	}

	Redo(st)

	assert.Equal(t, 1, st.HistoryIndex, "HistoryIndex should increment after Redo")
	assert.Equal(t, []rune("Hello, World!"), st.TextBuffer[0], "TextBuffer should have the change redone")
	assert.Equal(t, 13, st.CurrentCol, "CurrentCol should reflect the redone change")
}
