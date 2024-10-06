package editor

import (
	"testing"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/stretchr/testify/assert"
)

func TestStartSelection(t *testing.T) {
	st := &state.State{
		CurrentRow: 1,
		CurrentCol: 5,
	}

	startSelection(st)

	assert.Equal(t, 1, st.StartRow, "StartRow should match CurrentRow")
	assert.Equal(t, 5, st.StartCol, "StartCol should match CurrentCol")
	assert.True(t, st.SelectionActive, "Selection should be active")
}

func TestUpdateSelection(t *testing.T) {
	st := &state.State{
		CurrentRow: 2,
		CurrentCol: 3,
		StartRow:   1,
		StartCol:   5,
		SelectionActive: true,
	}

	updateSelection(st)

	assert.Equal(t, 2, st.EndRow, "EndRow should match CurrentRow")
	assert.Equal(t, 4, st.EndCol, "EndCol should match CurrentCol + 1")
}

func TestEndSelection(t *testing.T) {
	st := &state.State{
		SelectionActive: true,
	}

	endSelection(st)

	assert.False(t, st.SelectionActive, "Selection should be inactive")
}

func TestCopySelection(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello World"),
			[]rune("This is a test"),
		},
		StartRow: 0,
		StartCol: 6,
		EndRow:   0,
		EndCol:   11,
	}

	copySelection(st)

	expected := []rune("World")
	assert.Equal(t, expected, st.CopyBuffer, "CopyBuffer should contain 'World'")
}

func TestCutSelection(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello World"),
			[]rune("This is a test"),
		},
		StartRow: 0,
		StartCol: 6,
		EndRow:   0,
		EndCol:   11,
	}

	cutSelection(st)

	expectedCut := []rune("World")
	expectedBuffer := [][]rune{
		[]rune("Hello "),
		[]rune("This is a test"),
	}
	assert.Equal(t, expectedCut, st.CopyBuffer, "CopyBuffer should contain 'World'")
	assert.Equal(t, expectedBuffer, st.TextBuffer, "TextBuffer should have 'World' removed")
}

func TestPasteSelection(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello"),
		},
		CurrentRow: 0,
		CurrentCol: 5,
		CopyBuffer: []rune(" World"),
	}

	pasteSelection(st)

	expectedBuffer := [][]rune{
		[]rune("Hello World"),
	}
	assert.Equal(t, expectedBuffer, st.TextBuffer, "TextBuffer should contain 'Hello World'")
}

func TestDeleteSelection(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello World"),
			[]rune("This is a test"),
		},
		StartRow: 0,
		StartCol: 6,
		EndRow:   0,
		EndCol:   11,
	}

	deleteSelection(st)

	expectedBuffer := [][]rune{
		[]rune("Hello "),
		[]rune("This is a test"),
	}
	assert.Equal(t, expectedBuffer, st.TextBuffer, "TextBuffer should have 'World' removed")
}
