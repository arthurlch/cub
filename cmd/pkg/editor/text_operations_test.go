package editor_test

import (
	"testing"

	"github.com/arthurlch/cub/cmd/pkg/editor"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/stretchr/testify/assert"
)

func TestInsertRune(t *testing.T) {
	st := &state.State{
		Buffer: state.Buffer{TextBuffer: [][]rune{
			[]rune("Hello, World!"),
		}},
		Cursor: state.Cursor{CurrentRow: 0, CurrentCol: 7},
	}

	editor.InsertRune(st, 'X')

	assert.Equal(t, "Hello, XWorld!", string(st.TextBuffer[0]))
	assert.Equal(t, 8, st.CurrentCol)
	assert.True(t, st.Modified)
	assert.True(t, st.CanUndo())
}

func TestBackspace(t *testing.T) {
	st := &state.State{
		Buffer: state.Buffer{TextBuffer: [][]rune{
			[]rune("Hello, World!"),
		}},
		Cursor: state.Cursor{CurrentRow: 0, CurrentCol: 5},
	}

	editor.Backspace(st)

	assert.Equal(t, "Hell, World!", string(st.TextBuffer[0]))
	assert.Equal(t, 4, st.CurrentCol)
	assert.True(t, st.Modified)
	assert.True(t, st.CanUndo())
}

func TestBackspaceJoinsLines(t *testing.T) {
	st := &state.State{
		Buffer: state.Buffer{TextBuffer: [][]rune{
			[]rune("Hello"),
			[]rune("World!"),
		}},
		Cursor: state.Cursor{CurrentRow: 1, CurrentCol: 0},
	}

	editor.Backspace(st)

	assert.Equal(t, 1, len(st.TextBuffer))
	assert.Equal(t, "HelloWorld!", string(st.TextBuffer[0]))
	assert.Equal(t, 0, st.CurrentRow)
	assert.Equal(t, 5, st.CurrentCol)
	assert.True(t, st.Modified)
	assert.True(t, st.CanUndo())
}

func TestInsertNewLine(t *testing.T) {
	st := &state.State{
		Buffer: state.Buffer{TextBuffer: [][]rune{
			[]rune("Hello, World!"),
		}},
		Cursor: state.Cursor{CurrentRow: 0, CurrentCol: 7},
	}

	editor.InsertNewLine(st)

	assert.Equal(t, "Hello, ", string(st.TextBuffer[0]))
	assert.Equal(t, "World!", string(st.TextBuffer[1]))
	assert.Equal(t, 1, st.CurrentRow)
	assert.Equal(t, 0, st.CurrentCol)
	assert.True(t, st.Modified)
	assert.True(t, st.CanUndo())
}

func TestDeleteCurrentLine(t *testing.T) {
	st := &state.State{
		Buffer: state.Buffer{TextBuffer: [][]rune{
			[]rune("Hello, World!"),
			[]rune("This is a test."),
		}},
		Cursor: state.Cursor{CurrentRow: 0, CurrentCol: 0},
	}

	editor.DeleteCurrentLine(st)

	assert.Equal(t, "This is a test.", string(st.TextBuffer[0]))
	assert.Equal(t, 0, st.CurrentRow)
	assert.Equal(t, 0, st.CurrentCol)
	assert.True(t, st.Modified)
	assert.True(t, st.CanUndo())
}

func TestDeleteLastRemainingLine(t *testing.T) {
	st := &state.State{
		Buffer: state.Buffer{TextBuffer: [][]rune{
			[]rune("Only line."),
		}},
		Cursor: state.Cursor{CurrentRow: 0, CurrentCol: 5},
	}

	editor.DeleteCurrentLine(st)

	assert.Equal(t, 1, len(st.TextBuffer))
	assert.Equal(t, "", string(st.TextBuffer[0]))
	assert.Equal(t, 0, st.CurrentRow)
	assert.Equal(t, 0, st.CurrentCol)
	assert.True(t, st.Modified)
	assert.True(t, st.CanUndo())
}
