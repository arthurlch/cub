package state_test

import (
	"testing"
	"time"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/stretchr/testify/assert"
)

func TestActiveMessageNoneSet(t *testing.T) {
	s := &state.State{}
	_, ok := s.ActiveMessage()
	assert.False(t, ok, "no message should be active on a fresh state")
}

func TestActiveMessageShown(t *testing.T) {
	s := &state.State{}
	s.ShowMessage("File saved successfully.")

	msg, ok := s.ActiveMessage()
	assert.True(t, ok)
	assert.Equal(t, "File saved successfully.", msg)
}

func TestClampCursorEmptyBufferDoesNotPanic(t *testing.T) {
	s := &state.State{Cursor: state.Cursor{CurrentRow: 3, CurrentCol: 9}}
	assert.NotPanics(t, func() { s.ClampCursor() })
	assert.NotPanics(t, func() { s.MoveCursor(5, 5) })
	assert.Equal(t, 0, s.CurrentRow)
	assert.Equal(t, 0, s.CurrentCol)
}

func TestActiveMessageExpires(t *testing.T) {
	s := &state.State{}
	s.ShowMessage("File saved successfully.")
	s.MessageTimestamp = time.Now().Add(-time.Hour)

	_, ok := s.ActiveMessage()
	assert.False(t, ok, "message should expire after its TTL")
}
