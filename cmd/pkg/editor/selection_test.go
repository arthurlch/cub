package editor_test

import (
	"fmt"
	"testing"

	"github.com/arthurlch/cub/cmd/pkg/editor"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/stretchr/testify/assert"
)

func TestStartSelection(t *testing.T) {
	s := &state.State{
		Cursor: state.Cursor{CurrentRow: 1, CurrentCol: 3},
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
		Buffer: state.Buffer{TextBuffer: [][]rune{
			[]rune("Hello, World!"),
			[]rune("This is a test."),
			[]rune("Another line."),
		}},
		Cursor:    state.Cursor{CurrentRow: 2, CurrentCol: 5},
		Selection: state.Selection{StartRow: 1, StartCol: 3, SelectionActive: true},
	}
	editor.UpdateSelection(s)
	assert.Equal(t, 2, s.EndRow)
	assert.Equal(t, 5, s.EndCol)
	assert.True(t, s.SelectionActive)
}

func TestEndSelection(t *testing.T) {
	s := &state.State{
		Selection: state.Selection{SelectionActive: true},
	}
	editor.EndSelection(s)
	assert.False(t, s.SelectionActive)
}

func TestCopySelection(t *testing.T) {
	s := &state.State{
		Buffer: state.Buffer{TextBuffer: [][]rune{
			[]rune("Hello, World!"),
			[]rune("This is a test."),
		}},
		Selection: state.Selection{StartRow: 0, StartCol: 7, EndRow: 0, EndCol: 12},
	}
	editor.CopySelection(s)
	assert.Equal(t, []rune("World"), s.CopyBuffer)
}

func TestCutSelection(t *testing.T) {
	s := &state.State{
		Buffer: state.Buffer{TextBuffer: [][]rune{
			[]rune("Hello, World!"),
			[]rune("This is a test."),
		}},
		Selection: state.Selection{StartRow: 0, StartCol: 7, EndRow: 0, EndCol: 12},
	}
	editor.CutSelection(s)
	assert.Equal(t, []rune("World"), s.CopyBuffer)
	assert.Equal(t, []rune("Hello, !"), s.TextBuffer[0])
}

func TestPasteSelection(t *testing.T) {
	s := &state.State{
		Buffer: state.Buffer{TextBuffer: [][]rune{
			[]rune("Hello"),
			[]rune("This is a test."),
		}},
		Cursor:     state.Cursor{CurrentRow: 0, CurrentCol: 5},
		CopyBuffer: []rune(" World!"),
	}
	editor.PasteSelection(s)
	assert.Equal(t, [][]rune{
		[]rune("Hello World!"),
		[]rune("This is a test."),
	}, s.TextBuffer)
}

func TestPasteSelectionMidLine(t *testing.T) {
	s := &state.State{
		Buffer:     state.Buffer{TextBuffer: [][]rune{[]rune("abcdef")}},
		Cursor:     state.Cursor{CurrentRow: 0, CurrentCol: 2},
		CopyBuffer: []rune("XY"),
	}
	editor.PasteSelection(s)
	assert.Equal(t, "abXYcdef", string(s.TextBuffer[0]))
	assert.Equal(t, 4, s.CurrentCol)
}

func TestCopySelectionDoesNotModifyBuffer(t *testing.T) {
	s := &state.State{
		Buffer:    state.Buffer{TextBuffer: [][]rune{[]rune("Hello")}},
		Selection: state.Selection{StartRow: 0, StartCol: 0, EndRow: 0, EndCol: 6, SelectionActive: true},
	}
	editor.CopySelection(s)
	assert.Equal(t, "Hello", string(s.TextBuffer[0]), "copy must not modify the document")
	assert.Equal(t, "Hello", string(s.CopyBuffer))
}

func TestDeleteSelection(t *testing.T) {
	s := &state.State{
		Buffer: state.Buffer{TextBuffer: [][]rune{
			[]rune("Hello, World!"),
			[]rune("This is a test."),
		}},
		Selection: state.Selection{StartRow: 0, StartCol: 7, EndRow: 0, EndCol: 12},
	}
	editor.DeleteSelection(s)
	assert.Equal(t, [][]rune{
		[]rune("Hello, !"),
		[]rune("This is a test."),
	}, s.TextBuffer)
}

func TestSelectAll(t *testing.T) {
	s := &state.State{
		Buffer: state.Buffer{TextBuffer: [][]rune{
			[]rune("Hello, World!"),
			[]rune("This is a test."),
		}},
	}

	editor.SelectAll(s)

	fmt.Printf("StartRow: %d, StartCol: %d, EndRow: %d, EndCol: %d, LastLineLen: %d\n",
		s.StartRow, s.StartCol, s.EndRow, s.EndCol, len(s.TextBuffer[s.EndRow]))

	assert.Equal(t, 0, s.StartRow)
	assert.Equal(t, 0, s.StartCol)
	assert.Equal(t, 1, s.EndRow)
	assert.GreaterOrEqual(t, s.EndCol, 0, "EndCol should be within the valid range")
	assert.LessOrEqual(t, s.EndCol, len(s.TextBuffer[s.EndRow]), "EndCol should not exceed the last line's length")
	assert.True(t, s.SelectionActive)
}
