package editor_test

import (
	"testing"

	"github.com/arthurlch/cub/cmd/pkg/editor"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestUndoRedo(t *testing.T) {
	initialTextBuffer := [][]rune{
		[]rune("Hello"),
		[]rune("Cub"),
	}
	stateInstance := &state.State{
		TextBuffer: utils.DeepCopyTextBuffer(initialTextBuffer),
		CurrentRow: 0,
		CurrentCol: 5,
	}

	editedTextBuffer := [][]rune{
		[]rune("Hello,"),
		[]rune("Cub"),
	}
	stateInstance.UndoBuffer = append(stateInstance.UndoBuffer, state.UndoState{
		TextBuffer: utils.DeepCopyTextBuffer(stateInstance.TextBuffer),
		CurrentRow: stateInstance.CurrentRow,
		CurrentCol: stateInstance.CurrentCol,
	})
	stateInstance.TextBuffer = utils.DeepCopyTextBuffer(editedTextBuffer)
	stateInstance.CurrentRow = 0
	stateInstance.CurrentCol = 6

	editor.Undo(stateInstance)
	assert.Equal(t, initialTextBuffer, stateInstance.TextBuffer, "Undo should revert to the original text buffer")
	assert.Equal(t, 0, stateInstance.CurrentRow, "Undo should revert to the original row position")
	assert.Equal(t, 5, stateInstance.CurrentCol, "Undo should revert to the original column position")

	editor.Redo(stateInstance)
	assert.Equal(t, editedTextBuffer, stateInstance.TextBuffer, "Redo should restore the edited text buffer")
	assert.Equal(t, 0, stateInstance.CurrentRow, "Redo should restore the edited row position")
	assert.Equal(t, 6, stateInstance.CurrentCol, "Redo should restore the edited column position")
}

func TestUndoEmptyBuffer(t *testing.T) {
	stateInstance := &state.State{}

	editor.Undo(stateInstance)
	assert.Empty(t, stateInstance.RedoBuffer, "RedoBuffer should remain empty when Undo is called on an empty UndoBuffer")
}

func TestRedoEmptyBuffer(t *testing.T) {
	stateInstance := &state.State{}

	editor.Redo(stateInstance)
	assert.Empty(t, stateInstance.UndoBuffer, "UndoBuffer should remain empty when Redo is called on an empty RedoBuffer")
}