// cmd/pkg/editor/selection_test.go

package editor

import (
	"testing"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/stretchr/testify/assert"
)



func TestStartSelection(t *testing.T) {
	st := &state.State{
		CurrentRow:      2,
		CurrentCol:      5,
		StartRow:        0,
		StartCol:        0,
		EndRow:          0,
		EndCol:          0,
		SelectionActive: false,
	}

	startSelection(st)

	assert.Equal(t, 2, st.StartRow, "StartRow should be set to CurrentRow")
	assert.Equal(t, 5, st.StartCol, "StartCol should be set to CurrentCol")
	assert.Equal(t, 2, st.EndRow, "EndRow should be set to CurrentRow")
	assert.Equal(t, 5, st.EndCol, "EndCol should be set to CurrentCol")
	assert.True(t, st.SelectionActive, "SelectionActive should be true")
}

func TestUpdateSelection_SameLine(t *testing.T) {
	st := &state.State{
		StartRow:        1,
		StartCol:        3,
		EndRow:          1,
		EndCol:          3,
		CurrentRow:      1,
		CurrentCol:      3,
		SelectionActive: true,
	}

	updateSelection(st)

	assert.Equal(t, 1, st.EndRow, "EndRow should remain the same")
	assert.Equal(t, 4, st.EndCol, "EndCol should increment by 1")
}

func TestUpdateSelection_DifferentLine(t *testing.T) {
	st := &state.State{
		StartRow:        1,
		StartCol:        3,
		EndRow:          1,
		EndCol:          3,
		CurrentRow:      2,
		CurrentCol:      5,
		SelectionActive: true,
	}

	updateSelection(st)

	assert.Equal(t, 2, st.EndRow, "EndRow should update to CurrentRow")
	assert.Equal(t, 5, st.EndCol, "EndCol should update to CurrentCol")
}

func TestCopySelection_SingleLine(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello, World!"),
		},
		StartRow:        0,
		StartCol:        7,
		EndRow:          0,
		EndCol:          12,
		SelectionActive: true,
	}

	copySelection(st)

	expectedCopy := []rune("World")
	assert.Equal(t, expectedCopy, st.CopyBuffer, "CopyBuffer should contain 'World'")
}

func TestCopySelection_MultipleLines(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Line 1"),
			[]rune("Line 2"),
			[]rune("Line 3"),
		},
		StartRow:        0,
		StartCol:        5,
		EndRow:          2,
		EndCol:          4,
		SelectionActive: true,
	}

	copySelection(st)

	expectedCopy := []rune("1\nLine 2\nLine")
	assert.Equal(t, expectedCopy, st.CopyBuffer, "CopyBuffer should contain '1\\nLine 2\\nLine'")
}

func TestCutSelection_SingleLine(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello, World!"),
		},
		StartRow:        0,
		StartCol:        7,
		EndRow:          0,
		EndCol:          12,
		SelectionActive: true,
	}

	cutSelection(st)

	expectedCopy := []rune("World")
	expectedText := "Hello, !"
	assert.Equal(t, expectedCopy, st.CopyBuffer, "CopyBuffer should contain 'World'")
	assert.Equal(t, []rune(expectedText), st.TextBuffer[0], "TextBuffer should have 'World' removed")
	assert.Equal(t, 0, st.CurrentRow, "CurrentRow should remain the same")
	assert.Equal(t, 7, st.CurrentCol, "CurrentCol should remain the same")
	assert.True(t, st.Modified, "State should be marked as modified")
}

func TestPasteSelection(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello, !"),
		},
		CopyBuffer: []rune("World"),
		CurrentRow: 0,
		CurrentCol: 7,
	}

	pasteSelection(st)

	expectedText := "Hello, World!"
	assert.Equal(t, []rune(expectedText), st.TextBuffer[0], "TextBuffer should have 'World' pasted")
	assert.Equal(t, 12, st.CurrentCol, "Cursor should move to position 12")
	assert.True(t, st.Modified, "State should be marked as modified")
}

func TestDeleteSelection(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello, World!"),
		},
		StartRow:        0,
		StartCol:        7,
		EndRow:          0,
		EndCol:          12,
		SelectionActive: true,
	}

	deleteSelection(st)

	expectedText := "Hello, !"
	assert.Equal(t, []rune(expectedText), st.TextBuffer[0], "TextBuffer should have 'World' deleted")
	assert.Equal(t, 0, st.CurrentRow, "CurrentRow should remain the same")
	assert.Equal(t, 7, st.CurrentCol, "CurrentCol should remain the same")
	assert.True(t, st.Modified, "State should be marked as modified")
}
