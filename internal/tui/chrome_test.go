package tui

import (
	"strings"
	"testing"

	"github.com/arthurlch/cub/internal/editor"
	"github.com/arthurlch/cub/internal/filetree"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestLerpColorMidpoint(t *testing.T) {
	got := lerpColor(lipgloss.Color("#000000"), lipgloss.Color("#ffffff"), 0.5)
	assert.Equal(t, lipgloss.Color("#7f7f7f"), got)
}

func TestScrollbarHeightAndThumb(t *testing.T) {
	m := New(editor.New(), &filetree.Node{Name: "root"}, 26)
	out := m.scrollbar(10, 0, 100)
	assert.Equal(t, 10, len(strings.Split(out, "\n")), "one row per height unit")
	assert.Contains(t, out, "┃", "thumb shown when content overflows the viewport")
}

func TestScrollbarNoThumbWhenFits(t *testing.T) {
	m := New(editor.New(), &filetree.Node{Name: "root"}, 26)
	out := m.scrollbar(10, 0, 5)
	assert.NotContains(t, out, "┃")
}
