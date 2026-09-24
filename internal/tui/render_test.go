package tui

import (
	"strings"
	"testing"

	"github.com/arthurlch/cub/internal/editor"
	"github.com/arthurlch/cub/internal/filetree"
	"github.com/arthurlch/cub/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func init() { lipgloss.SetColorProfile(termenv.TrueColor) }

func sizedModel(t *testing.T, sidebar bool) Model {
	t.Helper()
	e := editor.New()
	e.SetText("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hi\")\n}")
	m := New(e, &filetree.Node{Name: "root"}, 26)
	tm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = tm.(Model)
	if sidebar {
		m.sbVisible, m.sbFocused, m.sbReveal = true, true, 1
		m.applySize()
	}
	return m
}

func TestViewFillsExactDimensions(t *testing.T) {
	m := sizedModel(t, false)
	view := m.View()
	assert.Equal(t, 30, lipgloss.Height(view), "view height must equal terminal height")
	for i, line := range strings.Split(view, "\n") {
		assert.Equalf(t, 100, lipgloss.Width(line), "row %d width", i)
	}
}

func TestViewWithSidebarFillsDimensions(t *testing.T) {
	m := sizedModel(t, true)
	view := m.View()
	assert.Equal(t, 30, lipgloss.Height(view))
	assert.Equal(t, 100, lipgloss.Width(view))
}

func TestViewRendersLightTheme(t *testing.T) {
	m := sizedModel(t, true)
	for i, th := range theme.Themes {
		if th.Name == "github" {
			m.themeIdx, m.theme, m.st = i, th, buildStyles(th)
			break
		}
	}
	view := m.View()
	assert.Equal(t, 30, lipgloss.Height(view))
	for i, line := range strings.Split(view, "\n") {
		assert.Equalf(t, 100, lipgloss.Width(line), "light theme row %d", i)
	}
}

func TestViewRightSidebarFillsDimensions(t *testing.T) {
	m := sizedModel(t, true)
	m.sbRight = true
	m.applySize()
	view := m.View()
	assert.Equal(t, 30, lipgloss.Height(view))
	for i, line := range strings.Split(view, "\n") {
		assert.Equalf(t, 100, lipgloss.Width(line), "right-sidebar row %d", i)
	}
}

func TestOverlaysRenderFullScreen(t *testing.T) {
	m := sizedModel(t, false)
	m.showHelp = true
	assert.Equal(t, 30, lipgloss.Height(m.View()))
	m.showHelp = false
	m.showTheme = true
	assert.Equal(t, 30, lipgloss.Height(m.View()))
}
