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

	es := &EditorState{State: st}
	keyEvent := termbox.Event{Ch: '!'}
	fileType := "go"

	es.InsertRunes(keyEvent, fileType)

	expectedBuffer := [][]rune{
		[]rune("Hello!"),
	}
	assert.Equal(t, expectedBuffer, st.TextBuffer, "TextBuffer should contain 'Hello!'")
	assert.Equal(t, 6, st.CurrentCol, "CurrentCol should be incremented after insertion")
	assert.True(t, st.Modified, "State should be modified")
}

func TestDeleteRune(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello!"),
		},
		CurrentRow: 0,
		CurrentCol: 6,
	}

	es := &EditorState{State: st}
	fileType := "go"

	es.DeleteRune(fileType)

	expectedBuffer := [][]rune{
		[]rune("Hello"),
	}
	assert.Equal(t, expectedBuffer, st.TextBuffer, "TextBuffer should contain 'Hello' after deletion")
	assert.Equal(t, 5, st.CurrentCol, "CurrentCol should be decremented after deletion")
	assert.True(t, st.Modified, "State should be modified")
}

func TestDeleteRuneAcrossLines(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello"),
			[]rune("World"),
		},
		CurrentRow: 1,
		CurrentCol: 0,
	}

	es := &EditorState{State: st}
	fileType := "go"

	es.DeleteRune(fileType)

	expectedBuffer := [][]rune{
		[]rune("HelloWorld"),
	}
	assert.Equal(t, expectedBuffer, st.TextBuffer, "TextBuffer should merge lines after deleting rune at start of line")
	assert.Equal(t, 5, st.CurrentCol, "CurrentCol should match the end of the previous line")
	assert.Equal(t, 0, st.CurrentRow, "CurrentRow should be decremented")
	assert.True(t, st.Modified, "State should be modified")
}

func TestInsertNewLine(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("Hello World"),
		},
		CurrentRow: 0,
		CurrentCol: 5,
	}

	es := &EditorState{State: st}
	fileType := "go"

	es.InsertNewLine(fileType) 

	expectedBuffer := [][]rune{
		[]rune("Hello"),
		[]rune(" World"),
	}
	assert.Equal(t, expectedBuffer, st.TextBuffer, "TextBuffer should split at CurrentCol with a new line inserted")
	assert.Equal(t, 1, st.CurrentRow, "CurrentRow should be incremented after newline insertion")
	assert.Equal(t, 0, st.CurrentCol, "CurrentCol should be reset to 0 after newline insertion")
	assert.True(t, st.Modified, "State should be modified")
}

func TestDeleteCurrentLine(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("First line"),
			[]rune("Second line"),
			[]rune("Third line"),
		},
		CurrentRow: 1,
		CurrentCol: 0,
	}

	deleteCurrentLine(st) 

	expectedBuffer := [][]rune{
		[]rune("First line"),
		[]rune("Third line"),
	}
	assert.Equal(t, expectedBuffer, st.TextBuffer, "TextBuffer should remove the current line")
	assert.Equal(t, 1, st.CurrentRow, "CurrentRow should be updated correctly")
	assert.Equal(t, 0, st.CurrentCol, "CurrentCol should be reset to 0 after deleting the line")
	assert.True(t, st.Modified, "State should be modified")
}

func TestDeleteCurrentLineWhenLastLine(t *testing.T) {
	st := &state.State{
		TextBuffer: [][]rune{
			[]rune("First line"),
		},
		CurrentRow: 0,
		CurrentCol: 5,
	}

	deleteCurrentLine(st) 

	expectedBuffer := [][]rune{
		{},
	}
	assert.Equal(t, expectedBuffer, st.TextBuffer, "TextBuffer should contain an empty line after deleting the last line")
	assert.Equal(t, 0, st.CurrentRow, "CurrentRow should be reset to 0")
	assert.Equal(t, 0, st.CurrentCol, "CurrentCol should be reset to 0")
	assert.True(t, st.Modified, "State should be modified")
}
