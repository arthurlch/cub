package tui

import (
	"testing"

	"github.com/arthurlch/cub/internal/editor"
	"github.com/arthurlch/cub/internal/filetree"
	"github.com/stretchr/testify/assert"
)

func TestSidebarFocusModel(t *testing.T) {
	m := &Model{e: editor.New(), tree: &filetree.Node{Name: "root"}}

	m.toggleSidebar()
	assert.True(t, m.sbVisible)
	assert.True(t, m.sbFocused, "opening focuses the sidebar")

	m.toggleSidebar()
	assert.True(t, m.sbVisible, "focus flip keeps it open")
	assert.False(t, m.sbFocused, "flips focus to the editor")

	m.toggleSidebar()
	assert.True(t, m.sbFocused, "flips focus back to the sidebar")

	m.hideSidebar()
	assert.False(t, m.sbVisible)
	assert.False(t, m.sbFocused)
	assert.Equal(t, 0.0, m.sbTarget, "hide animates toward closed")
}
