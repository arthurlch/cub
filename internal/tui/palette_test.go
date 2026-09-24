package tui

import (
	"testing"

	"github.com/arthurlch/cub/internal/editor"
	"github.com/arthurlch/cub/internal/filetree"
	"github.com/arthurlch/cub/internal/input"
	"github.com/stretchr/testify/assert"
)

func TestPaletteFilterAndOpen(t *testing.T) {
	e := editor.New()
	m := New(e, &filetree.Node{Name: "root"}, 26)
	m.showPalette = true
	m.paletteFiles = []string{"internal/theme/colors.go", "cmd/cub/main.go", "README.md"}
	m.recomputeMatches()
	assert.Equal(t, 3, len(m.paletteMatches))

	for _, r := range "main" {
		m.paletteKey(input.Event{Key: input.KeyRune, Rune: r})
	}
	assert.NotEmpty(t, m.paletteMatches)

	m.paletteKey(input.Event{Key: input.KeyEnter})
	assert.False(t, m.showPalette)
	assert.Equal(t, "cmd/cub/main.go", e.FileName())
}

func TestPaletteBackspaceOnEmpty(t *testing.T) {
	m := New(editor.New(), &filetree.Node{Name: "root"}, 26)
	m.showPalette = true
	assert.NotPanics(t, func() {
		m.paletteKey(input.Event{Key: input.KeyBackspace})
	})
}
