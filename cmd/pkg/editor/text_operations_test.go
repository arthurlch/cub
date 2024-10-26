package editor_test

import (
	"testing"

	"github.com/arthurlch/cub/cmd/pkg/editor"
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
)

func TestInsertRunes(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello, World!"),
		},
		CurrentRow: 0,
		CurrentCol: 7,
	}
	es := &editor.EditorState{State: st}

	// Simulate inserting a rune
	keyEvent := termbox.Event{Ch: 'X'}
	es.InsertRunes(keyEvent, "txt")

	assert.Equal(t, "Hello, XWorld!", string(st.TextBuffer[0]))
	assert.Equal(t, 8, st.CurrentCol)
	assert.True(t, st.Modified)
	assert.Len(t, st.UndoBuffer, 1)
}

func TestDeleteRune(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello, World!"),
		},
		CurrentRow: 0,
		CurrentCol: 5,
	}
	es := &editor.EditorState{State: st}

	// Simulate deleting a rune
	es.DeleteRune("txt")

	assert.Equal(t, "Hell, World!", string(st.TextBuffer[0]))
	assert.Equal(t, 4, st.CurrentCol)
	assert.True(t, st.Modified)
	assert.Len(t, st.UndoBuffer, 1)
}

func TestDeleteRuneAtLineStart(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello"),
			[]rune("World!"),
		},
		CurrentRow: 1,
		CurrentCol: 0,
	}
	es := &editor.EditorState{State: st}

	// Attempt to delete at the start of the line (no-op)
	es.DeleteRune("txt")

	assert.Equal(t, "Hello", string(st.TextBuffer[0]))
	assert.Equal(t, "World!", string(st.TextBuffer[1]))
	assert.Equal(t, 0, st.CurrentCol)
	assert.False(t, st.Modified)
}

func TestInsertNewLine(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello, World!"),
		},
		CurrentRow: 0,
		CurrentCol: 7,
	}
	es := &editor.EditorState{State: st}

	// Simulate inserting a new line
	es.InsertNewLine("txt")

	assert.Equal(t, "Hello, ", string(st.TextBuffer[0]))
	assert.Equal(t, "World!", string(st.TextBuffer[1]))
	assert.Equal(t, 1, st.CurrentRow)
	assert.Equal(t, 0, st.CurrentCol)
	assert.True(t, st.Modified)
	assert.Len(t, st.UndoBuffer, 1)
}

func TestDeleteCurrentLine(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello, World!"),
			[]rune("This is a test."),
		},
		CurrentRow: 0,
		CurrentCol: 0,
	}
	
	editor.DeleteCurrentLine(st)

	assert.Equal(t, "This is a test.", string(st.TextBuffer[0]))
	assert.Equal(t, 0, st.CurrentRow)
	assert.Equal(t, 0, st.CurrentCol)
	assert.True(t, st.Modified)
	assert.Len(t, st.UndoBuffer, 1)
}

func TestDeleteLastRemainingLine(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Only line."),
		},
		CurrentRow: 0,
		CurrentCol: 5,
	}

	editor.DeleteCurrentLine(st)

	assert.Equal(t, 1, len(st.TextBuffer))
	assert.Equal(t, "", string(st.TextBuffer[0]))
	assert.Equal(t, 0, st.CurrentRow)
	assert.Equal(t, 0, st.CurrentCol)
	assert.True(t, st.Modified)
	assert.Len(t, st.UndoBuffer, 1)
}
