// cmd/pkg/editor/text_operations_test.go

package editor

import (
	"testing"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
)

func TestInsertRunes(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello"),
		},
		CurrentRow: 0,
		CurrentCol: 5,
	}

	es := &EditorState{
		State: st,
	}

	keyEvent := termbox.Event{
		Key: termbox.Key(0),
		Ch:  '!',
	}

	es.InsertRunes(keyEvent)

	expectedText := "Hello!"
	assert.Equal(t, []rune(expectedText), st.TextBuffer[0], "TextBuffer should have 'Hello!' after insertion")
	assert.Equal(t, 6, st.CurrentCol, "Cursor should move to position 6")
}

func TestDeleteRune(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello!"),
		},
		CurrentRow: 0,
		CurrentCol: 6,
	}

	es := &EditorState{
		State: st,
	}

	es.DeleteRune()

	expectedText := "Hello"
	assert.Equal(t, []rune(expectedText), st.TextBuffer[0], "TextBuffer should have 'Hello' after deletion")
	assert.Equal(t, 5, st.CurrentCol, "Cursor should move to position 5")
}

func TestInsertNewLine(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello World"),
		},
		CurrentRow: 0,
		CurrentCol: 5,
	}

	es := &EditorState{
		State: st,
	}

	es.InsertNewLine()

	expectedTextLine1 := "Hello"
	expectedTextLine2 := " World"
	assert.Equal(t, 2, len(st.TextBuffer), "TextBuffer should have two lines after inserting new line")
	assert.Equal(t, []rune(expectedTextLine1), st.TextBuffer[0], "First line should be 'Hello'")
	assert.Equal(t, []rune(expectedTextLine2), st.TextBuffer[1], "Second line should be ' World'")
	assert.Equal(t, 1, st.CurrentRow, "CurrentRow should move to the new line")
	assert.Equal(t, 0, st.CurrentCol, "CurrentCol should reset to 0")
}

func TestDeleteCurrentLine(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Line 1"),
			[]rune("Line 2"),
			[]rune("Line 3"),
		},
		CurrentRow: 1,
		CurrentCol: 5,
		HistoryIndex: 2,
		ChangeHistory: []state.Change{
		},
	}

	deleteCurrentLine(st)

	expectedText := [][]rune{
		[]rune("Line 1"),
		[]rune("Line 3"),
	}
	assert.Equal(t, expectedText, st.TextBuffer, "TextBuffer should have 'Line 2' removed")
	assert.Equal(t, 1, st.CurrentRow, "CurrentRow should remain at the same index")
	assert.Equal(t, 5, st.CurrentCol, "CurrentCol should remain the same")
	assert.Equal(t, 3, st.HistoryIndex, "HistoryIndex should increment")
}
