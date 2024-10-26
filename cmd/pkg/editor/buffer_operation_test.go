package editor_test

import (
	"testing"

	"github.com/arthurlch/cub/cmd/pkg/editor"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/stretchr/testify/assert"
)

func TestUndo(t *testing.T) {
	s := &state.State{
		TextBuffer: [][]rune{{'H', 'e', 'l', 'l', 'o'}, {'W', 'o', 'r', 'l', 'd'}},
		CurrentRow: 1,
		CurrentCol: 5,
		UndoBuffer: []state.UndoState{
			{
				TextBuffer: [][]rune{{'H', 'i'}},
				CurrentRow: 0,
				CurrentCol: 2,
			},
		},
		RedoBuffer: []state.UndoState{},
		Modified:   false,
	}

	editor.Undo(s)

	assert.Equal(t, [][]rune{{'H', 'i'}}, s.TextBuffer, "Undo should revert text buffer to previous state")
	assert.Equal(t, 0, s.CurrentRow, "Undo should revert to previous row")
	assert.Equal(t, 2, s.CurrentCol, "Undo should revert to previous column")
	assert.True(t, s.Modified, "Undo operation should set Modified to true")
	assert.Equal(t, 1, len(s.RedoBuffer), "Redo buffer should contain the previous state after Undo")
}

func TestRedo(t *testing.T) {
	s := &state.State{
		TextBuffer: [][]rune{{'H', 'i'}},
		CurrentRow: 0,
		CurrentCol: 2,
		UndoBuffer: []state.UndoState{
			{
				TextBuffer: [][]rune{{'H', 'i'}},
				CurrentRow: 0,
				CurrentCol: 2,
			},
		},
		RedoBuffer: []state.UndoState{
			{
				TextBuffer: [][]rune{{'H', 'e', 'l', 'l', 'o'}, {'W', 'o', 'r', 'l', 'd'}},
				CurrentRow: 1,
				CurrentCol: 5,
			},
		},
		Modified: false,
	}


	editor.Redo(s)

	assert.Equal(t, [][]rune{{'H', 'e', 'l', 'l', 'o'}, {'W', 'o', 'r', 'l', 'd'}}, s.TextBuffer, "Redo should revert text buffer to previous undone state")
	assert.Equal(t, 1, s.CurrentRow, "Redo should revert to previous undone row")
	assert.Equal(t, 5, s.CurrentCol, "Redo should revert to previous undone column")
	assert.True(t, s.Modified, "Redo operation should set Modified to true")
	assert.Equal(t, 1, len(s.UndoBuffer), "Undo buffer should contain the state after Redo")
}

func TestUndoWithoutUndoBuffer(t *testing.T) {
	s := &state.State{
		TextBuffer: [][]rune{{'H', 'e', 'l', 'l', 'o'}},
		CurrentRow: 0,
		CurrentCol: 5,
		UndoBuffer: []state.UndoState{},
		RedoBuffer: []state.UndoState{},
		Modified:   false,
	}

	editor.Undo(s)

	assert.Equal(t, [][]rune{{'H', 'e', 'l', 'l', 'o'}}, s.TextBuffer, "TextBuffer should remain unchanged if UndoBuffer is empty")
	assert.Equal(t, 0, s.CurrentRow, "Row should remain unchanged if UndoBuffer is empty")
	assert.Equal(t, 5, s.CurrentCol, "Column should remain unchanged if UndoBuffer is empty")
	assert.False(t, s.Modified, "Modified should remain false if UndoBuffer is empty")
}

func TestRedoWithoutRedoBuffer(t *testing.T) {
	s := &state.State{
		TextBuffer: [][]rune{{'H', 'i'}},
		CurrentRow: 0,
		CurrentCol: 2,
		UndoBuffer: []state.UndoState{
			{
				TextBuffer: [][]rune{{'H', 'i'}},
				CurrentRow: 0,
				CurrentCol: 2,
			},
		},
		RedoBuffer: []state.UndoState{},
		Modified:   false,
	}

	editor.Redo(s)

	assert.Equal(t, [][]rune{{'H', 'i'}}, s.TextBuffer, "TextBuffer should remain unchanged if RedoBuffer is empty")
	assert.Equal(t, 0, s.CurrentRow, "Row should remain unchanged if RedoBuffer is empty")
	assert.Equal(t, 2, s.CurrentCol, "Column should remain unchanged if RedoBuffer is empty")
	assert.False(t, s.Modified, "Modified should remain false if RedoBuffer is empty")
}
