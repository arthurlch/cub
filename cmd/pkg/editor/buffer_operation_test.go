package editor_test

import (
	"log"
	"os"
	"testing"

	"github.com/arthurlch/cub/cmd/pkg/editor"
	"github.com/arthurlch/cub/cmd/pkg/logging"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/stretchr/testify/assert"
)

func init() {
	logging.Logger = log.New(os.Stdout, "TEST: ", log.LstdFlags)
}

func TestUndo(t *testing.T) {
	s := &state.State{
		Buffer: state.Buffer{TextBuffer: [][]rune{{'H', 'i'}}},
		Cursor: state.Cursor{CurrentRow: 0, CurrentCol: 2},
	}

	s.Do("edit", func() {
		s.TextBuffer = [][]rune{{'H', 'e', 'l', 'l', 'o'}, {'W', 'o', 'r', 'l', 'd'}}
		s.CurrentRow = 1
		s.CurrentCol = 5
	})
	assert.True(t, s.CanUndo())

	editor.Undo(s)

	assert.Equal(t, [][]rune{{'H', 'i'}}, s.TextBuffer, "Undo should revert text buffer to previous state")
	assert.Equal(t, 0, s.CurrentRow, "Undo should revert to previous row")
	assert.Equal(t, 2, s.CurrentCol, "Undo should revert to previous column")
	assert.True(t, s.Modified, "Undo operation should set Modified to true")
	assert.True(t, s.CanRedo(), "Redo stack should contain the undone command")
}

func TestRedo(t *testing.T) {
	s := &state.State{
		Buffer: state.Buffer{TextBuffer: [][]rune{{'H', 'i'}}},
		Cursor: state.Cursor{CurrentRow: 0, CurrentCol: 2},
	}
	s.Do("edit", func() {
		s.TextBuffer = [][]rune{{'H', 'e', 'l', 'l', 'o'}, {'W', 'o', 'r', 'l', 'd'}}
		s.CurrentRow = 1
		s.CurrentCol = 5
	})
	editor.Undo(s)

	editor.Redo(s)

	assert.Equal(t, [][]rune{{'H', 'e', 'l', 'l', 'o'}, {'W', 'o', 'r', 'l', 'd'}}, s.TextBuffer, "Redo should re-apply the undone edit")
	assert.Equal(t, 1, s.CurrentRow, "Redo should restore the row")
	assert.Equal(t, 5, s.CurrentCol, "Redo should restore the column")
	assert.True(t, s.Modified, "Redo operation should set Modified to true")
	assert.True(t, s.CanUndo(), "Undo stack should contain the redone command")
}

func TestUndoWithoutHistory(t *testing.T) {
	s := &state.State{
		Buffer: state.Buffer{TextBuffer: [][]rune{{'H', 'e', 'l', 'l', 'o'}}},
		Cursor: state.Cursor{CurrentRow: 0, CurrentCol: 5},
	}

	editor.Undo(s)

	assert.Equal(t, [][]rune{{'H', 'e', 'l', 'l', 'o'}}, s.TextBuffer, "TextBuffer should remain unchanged if there is nothing to undo")
	assert.Equal(t, 0, s.CurrentRow, "Row should remain unchanged if there is nothing to undo")
	assert.Equal(t, 5, s.CurrentCol, "Column should remain unchanged if there is nothing to undo")
	assert.False(t, s.Modified, "Modified should remain false if there is nothing to undo")
}

func TestRedoWithoutRedoStack(t *testing.T) {
	s := &state.State{
		Buffer: state.Buffer{TextBuffer: [][]rune{{'H', 'i'}}},
		Cursor: state.Cursor{CurrentRow: 0, CurrentCol: 2},
	}
	s.Do("edit", func() { s.TextBuffer = [][]rune{{'H', 'i', '!'}} })

	editor.Redo(s)

	assert.Equal(t, [][]rune{{'H', 'i', '!'}}, s.TextBuffer, "Redo with an empty redo stack should be a no-op")
	assert.False(t, s.Modified, "Modified should remain false if there is nothing to redo")
}
