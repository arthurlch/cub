package tui

import (
	"testing"

	"github.com/arthurlch/cub/internal/editor"
	"github.com/arthurlch/cub/internal/input"
	"github.com/arthurlch/cub/internal/theme"
	"github.com/stretchr/testify/assert"
)

func TestThemePickerApplyAndCancel(t *testing.T) {
	m := &Model{e: editor.New(), theme: theme.Default(), themeIdx: 0}
	m.st = buildStyles(m.theme)
	m.showTheme = true
	m.themeSel = 0

	m.themeKey(input.Event{Key: input.KeyDown})
	assert.Equal(t, 1, m.themeSel)
	assert.Equal(t, theme.Themes[1].Name, m.theme.Name, "moving previews the theme live")

	m.themeKey(input.Event{Key: input.KeyEnter})
	assert.False(t, m.showTheme)
	assert.Equal(t, 1, m.themeIdx)

	m.showTheme = true
	m.themeSel = 1
	m.themeKey(input.Event{Key: input.KeyDown})
	assert.Equal(t, theme.Themes[2].Name, m.theme.Name)
	m.themeKey(input.Event{Key: input.KeyEsc})
	assert.False(t, m.showTheme)
	assert.Equal(t, theme.Themes[1].Name, m.theme.Name, "Esc reverts to the committed theme")
}
