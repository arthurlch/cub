package tui

import (
	"testing"

	"github.com/arthurlch/cub/internal/input"
	"github.com/stretchr/testify/assert"
)

func TestSidebarResizeGrowShrink(t *testing.T) {
	m := &Model{sbWidth: 26, width: 120, sbVisible: true, sbFocused: true}

	m.sidebarKey(input.Event{Key: input.KeyRune, Rune: '>'})
	assert.Equal(t, 28, m.sbWidth)

	m.sidebarKey(input.Event{Key: input.KeyRune, Rune: '<'})
	assert.Equal(t, 26, m.sbWidth)
}

func TestSidebarResizeClamps(t *testing.T) {
	m := &Model{sbWidth: 12, width: 120, sbVisible: true, sbFocused: true}
	for i := 0; i < 10; i++ {
		m.sidebarKey(input.Event{Key: input.KeyRune, Rune: '<'})
	}
	assert.Equal(t, 12, m.sbWidth, "should not shrink below the minimum")

	m = &Model{sbWidth: 12, width: 40, sbVisible: true, sbFocused: true}
	for i := 0; i < 40; i++ {
		m.sidebarKey(input.Event{Key: input.KeyRune, Rune: '>'})
	}
	assert.LessOrEqual(t, m.sbWidth, 20, "should leave room for the editor")
}
