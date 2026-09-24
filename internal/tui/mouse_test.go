package tui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arthurlch/cub/internal/editor"
	"github.com/arthurlch/cub/internal/filetree"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestSidebarRightSideMouse(t *testing.T) {
	m := &Model{e: editor.New(), tree: &filetree.Node{Name: "root"}, sbWidth: 26, width: 120, height: 40, sbVisible: true, sbReveal: 1, sbRight: true}

	assert.Equal(t, 0, m.editorOriginX(), "editor sits at the left edge when the sidebar is on the right")
	assert.Equal(t, 94, m.sidebarBorderX())
	assert.True(t, m.inSidebarX(110))
	assert.False(t, m.inSidebarX(10))

	m.handleMouse(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonLeft, X: 94, Y: 5})
	assert.True(t, m.resizing)
	m.handleMouse(tea.MouseMsg{Action: tea.MouseActionMotion, X: 100, Y: 5})
	assert.Equal(t, 20, m.sbWidth, "dragging left from the right edge shrinks the sidebar")
}

func TestClickSidebarSelectsCorrectNode(t *testing.T) {
	dir := t.TempDir()
	for _, n := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(dir, n), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	tree, err := filetree.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	m := &Model{e: editor.New(), tree: tree, width: 100, height: 40, sbVisible: true, sbFocused: true, sbReveal: 1}

	// so here the first node sits at screen row 4 (title + top-pad + header + pad).
	// the easiest way assert this is to click the first two rows of the sidebar and check that the correct node is selected.
	m.clickSidebar(4)
	assert.Equal(t, 0, m.sbSelected)
	assert.Contains(t, m.e.FileName(), "a.txt")

	m.clickSidebar(5)
	assert.Equal(t, 1, m.sbSelected)
	assert.Contains(t, m.e.FileName(), "b.txt")
}

func TestMouseDragResizesSidebar(t *testing.T) {
	m := &Model{e: editor.New(), sbWidth: 26, width: 120, height: 40, sbVisible: true, sbReveal: 1}

	m.handleMouse(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonLeft, X: m.sbWidth - 1, Y: 5})
	assert.True(t, m.resizing)

	m.handleMouse(tea.MouseMsg{Action: tea.MouseActionMotion, X: 40, Y: 5})
	assert.Equal(t, 41, m.sbWidth)

	m.handleMouse(tea.MouseMsg{Action: tea.MouseActionRelease, X: 40, Y: 5})
	assert.False(t, m.resizing)
}

func TestMouseClickPlacesCursor(t *testing.T) {
	e := editor.New()
	e.SetText("hello world\nsecond line here")
	m := &Model{e: e, sbWidth: 26, width: 120, height: 40}
	m.applySize()

	m.handleMouse(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonLeft, X: gutterWidth + 6, Y: 2})
	r, c := e.Cursor()
	assert.Equal(t, 1, r)
	assert.Equal(t, 6, c)
}

func TestWheelScrollsEditor(t *testing.T) {
	e := editor.New()
	e.SetText("l0\nl1\nl2\nl3\nl4\nl5\nl6")
	m := &Model{e: e, sbWidth: 26, width: 120, height: 40}
	m.applySize()

	m.handleMouse(tea.MouseMsg{Button: tea.MouseButtonWheelDown, X: 40, Y: 5})
	r, _ := e.Cursor()
	assert.Equal(t, 3, r)
}
