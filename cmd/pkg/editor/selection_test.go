package editor_test

import (
	"testing"

	"github.com/arthurlch/cub/cmd/pkg/editor"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/stretchr/testify/assert"
)

func TestStartSelection(t *testing.T) {
	s := &state.State{
		CurrentRow:      1,
		CurrentCol:      3,
		SelectionActive: false,
	}
	editor.StartSelection(s)
	assert.Equal(t, 1, s.StartRow)
	assert.Equal(t, 3, s.StartCol)
	assert.Equal(t, 1, s.EndRow)
	assert.Equal(t, 3, s.EndCol)
	assert.True(t, s.SelectionActive)
}

func TestUpdateSelection(t *testing.T) {
	s := &state.State{
		CurrentRow: 2,
		CurrentCol: 5,
		StartRow:   1,
		StartCol:   3,
	}
	editor.UpdateSelection(s)
	assert.Equal(t, 2, s.EndRow)
	assert.Equal(t, 5, s.EndCol)
}

func TestEndSelection(t *testing.T) {
	s := &state.State{
		SelectionActive: true,
	}
	editor.EndSelection(s)
	assert.False(t, s.SelectionActive)
}

func TestCopySelection(t *testing.T) {
	s := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello, World!"),
			[]rune("This is a test."),
		},
		StartRow: 0,
		StartCol: 7,
		EndRow:   0,
		EndCol:   12,
	}
	editor.CopySelection(s)
	assert.Equal(t, []rune("World"), s.CopyBuffer)
}

func TestCutSelection(t *testing.T) {
	s := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello, World!"),
			[]rune("This is a test."),
		},
		StartRow: 0,
		StartCol: 7,
		EndRow:   0,
		EndCol:   12,
	}
	editor.CutSelection(s)
	assert.Equal(t, []rune("World"), s.CopyBuffer)
	assert.Equal(t, []rune("Hello, !"), s.TextBuffer[0])
}

func TestPasteSelection(t *testing.T) {
	s := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello"),
			[]rune("This is a test."),
		},
		CurrentRow: 0,
		CurrentCol: 5,
		CopyBuffer: []rune(" World!"),
	}
	editor.PasteSelection(s)
	assert.Equal(t, [][]rune{
		[]rune("Hello World!"),
		[]rune("This is a test."),
	}, s.TextBuffer)
}

func TestDeleteSelection(t *testing.T) {
	s := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello, World!"),
			[]rune("This is a test."),
		},
		StartRow: 0,
		StartCol: 7,
		EndRow:   0,
		EndCol:   12,
	}
	editor.DeleteSelection(s)
	assert.Equal(t, [][]rune{
		[]rune("Hello, !"),
		[]rune("This is a test."),
	}, s.TextBuffer)
}

func TestSelectAll(t *testing.T) {
	s := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello, World!"),
			[]rune("This is a test."),
		},
	}
	editor.SelectAll(s)
	assert.Equal(t, 0, s.StartRow)
	assert.Equal(t, 0, s.StartCol)
	assert.Equal(t, 1, s.EndRow)
	assert.Equal(t, len(s.TextBuffer[1]), s.EndCol)
	assert.True(t, s.SelectionActive)
}
