package tui_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arthurlch/cub/internal/editor"
	"github.com/arthurlch/cub/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestModelEditsAndSaves(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(path, []byte("HELLO"), 0o644); err != nil {
		t.Fatal(err)
	}

	e := editor.New()
	if err := e.Open(path); err != nil {
		t.Fatal(err)
	}

	var m tea.Model = tui.New(e, nil, 26)
	m, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("i")})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("X")})
	m.Update(tea.KeyMsg{Type: tea.KeyCtrlS})

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "XHELLO", string(data))
}

func TestModelQuits(t *testing.T) {
	e := editor.New()
	var m tea.Model = tui.New(e, nil, 26)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlQ})
	assert.NotNil(t, cmd)
}

func TestViewRendersWithoutPanic(t *testing.T) {
	e := editor.New()
	e.SetText("package main\n\nfunc main() {}")
	var m tea.Model = tui.New(e, nil, 26)
	m, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	assert.NotEmpty(t, m.View())
}
