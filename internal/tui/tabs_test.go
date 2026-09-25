package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arthurlch/cub/internal/editor"
	"github.com/arthurlch/cub/internal/filetree"
	"github.com/arthurlch/cub/internal/syntax"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
)

func tmpFile(t *testing.T, dir, name, body string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestOpenInTabReusesScratchThenAddsTabs(t *testing.T) {
	dir := t.TempDir()
	a := tmpFile(t, dir, "a.txt", "aaa")
	b := tmpFile(t, dir, "b.txt", "bbb")

	m := New(editor.New(), &filetree.Node{Name: "root"}, 26)

	m.openInTab(a)
	assert.Equal(t, 1, len(m.docs), "first open reuses the empty scratch tab")
	assert.Equal(t, a, m.e.FileName())

	m.openInTab(b)
	assert.Equal(t, 2, len(m.docs), "second open adds a tab")
	assert.Equal(t, 1, m.active)
	assert.Equal(t, b, m.e.FileName())

	m.openInTab(a)
	assert.Equal(t, 2, len(m.docs), "reopening focuses the existing tab")
	assert.Equal(t, 0, m.active)
}

func TestTabSwitchAndClose(t *testing.T) {
	dir := t.TempDir()
	m := New(editor.New(), &filetree.Node{Name: "root"}, 26)
	m.openInTab(tmpFile(t, dir, "a.txt", "a"))
	m.openInTab(tmpFile(t, dir, "b.txt", "b"))
	m.openInTab(tmpFile(t, dir, "c.txt", "c"))
	assert.Equal(t, 3, len(m.docs))

	m.nextTab()
	assert.Equal(t, 0, m.active, "next wraps around")
	m.prevTab()
	assert.Equal(t, 2, m.active, "prev wraps around")

	m.setActive(1)
	m.closeTab()
	assert.Equal(t, 2, len(m.docs))

	m.closeTab()
	m.closeTab()
	assert.Equal(t, 1, len(m.docs), "closing the last tab leaves a fresh scratch buffer")
	assert.Equal(t, "", m.e.FileName())
}

func TestNewTabAddsAndActivates(t *testing.T) {
	m := New(editor.New(), &filetree.Node{Name: "root"}, 26)
	m.newTab()
	assert.Equal(t, 2, len(m.docs))
	assert.Equal(t, 1, m.active)
	assert.Equal(t, "", m.e.FileName(), "new tab is an empty scratch buffer")
}

func TestCloseGuardsUnsavedChanges(t *testing.T) {
	dir := t.TempDir()
	m := New(editor.New(), &filetree.Node{Name: "root"}, 26)
	m.openInTab(tmpFile(t, dir, "a.txt", "a"))
	m.openInTab(tmpFile(t, dir, "b.txt", "b"))
	m.e.InsertRune('x') // active tab now modified

	m.requestClose()
	assert.Equal(t, 2, len(m.docs), "first close on a dirty tab is a no-op warning")
	m.requestClose()
	assert.Equal(t, 1, len(m.docs), "second close discards")
}

func TestMouseClosesTabViaCloseGlyph(t *testing.T) {
	dir := t.TempDir()
	m := New(editor.New(), &filetree.Node{Name: "root"}, 26)
	tm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = tm.(Model)
	m.openInTab(tmpFile(t, dir, "a.txt", "a"))
	m.openInTab(tmpFile(t, dir, "b.txt", "b"))
	assert.Equal(t, 2, len(m.docs))

	boxes := m.tabLayout(m.width)
	first := boxes[0]
	m.clickTab(first.x + first.w - 1) // the ✕ region
	assert.Equal(t, 1, len(m.docs), "clicking the close glyph closes that tab")
}

func TestMouseCloseModifiedTabTakesTwoClicks(t *testing.T) {
	dir := t.TempDir()
	m := New(editor.New(), &filetree.Node{Name: "root"}, 26)
	tm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = tm.(Model)
	m.openInTab(tmpFile(t, dir, "a.txt", "a"))
	m.openInTab(tmpFile(t, dir, "b.txt", "b"))
	m.e.InsertRune('x') // active tab (b) now modified
	assert.Equal(t, 2, len(m.docs))

	boxes := m.tabLayout(m.width)
	closeX := boxes[1].x + boxes[1].w - 1

	m.clickTab(closeX)
	assert.Equal(t, 2, len(m.docs), "first click on a dirty tab's ✕ only warns")
	m.clickTab(closeX)
	assert.Equal(t, 1, len(m.docs), "second click on the ✕ discards and closes")
}

func TestMouseClickTabBodySwitchesWithoutClosing(t *testing.T) {
	dir := t.TempDir()
	m := New(editor.New(), &filetree.Node{Name: "root"}, 26)
	tm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = tm.(Model)
	m.openInTab(tmpFile(t, dir, "a.txt", "a"))
	m.openInTab(tmpFile(t, dir, "b.txt", "b"))

	boxes := m.tabLayout(m.width)
	first := boxes[0]
	m.clickTab(first.x + 1) // on the label, not the glyph
	assert.Equal(t, 2, len(m.docs), "clicking the tab body only switches")
	assert.Equal(t, 0, m.active)
}

func TestTabStopRendering(t *testing.T) {
	e := editor.New()
	e.SetText("\tx")
	m := New(e, &filetree.Node{Name: "root"}, 26)
	out := ansi.Strip(m.renderLine(0, syntax.GetLexer(""), 20, 0, false))
	runes := []rune(out)
	assert.Equal(t, 20, len(runes), "row is exactly width cells")
	assert.Equal(t, tabStop, strings.IndexRune(out, 'x'), "leading tab expands to the tab stop")
}
