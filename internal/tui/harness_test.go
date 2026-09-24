package tui

import (
	"strings"
	"testing"

	"github.com/arthurlch/cub/internal/editor"
	"github.com/arthurlch/cub/internal/filetree"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

type app struct {
	t *testing.T
	m Model
}

func newApp(t *testing.T, text string) *app {
	t.Helper()
	e := editor.New()
	e.SetText(text)
	m := New(e, &filetree.Node{Name: "root"}, 26)
	tm, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	return &app{t: t, m: tm.(Model)}
}

func (a *app) send(msg tea.Msg) *app {
	a.t.Helper()
	tm, _ := a.m.Update(msg)
	a.m = tm.(Model)
	return a
}

func (a *app) press(t tea.KeyType) *app {
	return a.send(tea.KeyMsg{Type: t})
}

func (a *app) typ(s string) *app {
	for _, r := range s {
		a.send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	return a
}

func (a *app) text() string { return a.m.e.Buffer().String() }

func (a *app) screen() string { return ansi.Strip(a.m.View()) }

func (a *app) shows(sub string) bool { return strings.Contains(a.screen(), sub) }
