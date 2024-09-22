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
				Col:  13,
				Text: []rune("!"), 
			},
		},
		CurrentRow: 0,
		CurrentCol: 13,
	}

	Undo(st)

	assert.Equal(t, 0, st.HistoryIndex, "HistoryIndex should decrement after Undo")
	assert.Equal(t, []rune("Hello, World"), st.TextBuffer[0], "TextBuffer should have the last change undone")
	assert.Equal(t, 0, st.CurrentRow, "CurrentRow should remain the same after Undo")
	assert.Equal(t, 12, st.CurrentCol, "CurrentCol should point to the position before the undo action")
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
	assert.Equal(t, 0, st.CurrentRow, "CurrentRow should remain unchanged after Redo")
	assert.Equal(t, 13, st.CurrentCol, "CurrentCol should reflect the redone change")
}

func TestInvertChange(t *testing.T) {
	change := state.Change{
		Type:    state.Insert,
		Row:     1,
		Col:     5,
		Text:    []rune("test"),
		PrevRow: 0,
		PrevCol: 0,
	}

	inverted := invertChange(change)

	assert.Equal(t, state.Delete, inverted.Type, "Inverted change should be Delete when original was Insert")
	assert.Equal(t, change.Row, inverted.Row, "Row should remain unchanged")
	assert.Equal(t, change.Col, inverted.Col, "Col should remain unchanged")
	assert.Equal(t, change.Text, inverted.Text, "Text should remain unchanged")
	assert.Equal(t, change.PrevRow, inverted.PrevRow, "PrevRow should remain unchanged")
	assert.Equal(t, change.PrevCol, inverted.PrevCol, "PrevCol should remain unchanged")
}

func TestInvertChangeType(t *testing.T) {
	insertType := invertChangeType(state.Insert)
	assert.Equal(t, state.Delete, insertType, "Insert should invert to Delete")

	deleteType := invertChangeType(state.Delete)
	assert.Equal(t, state.Insert, deleteType, "Delete should invert to Insert")
}
