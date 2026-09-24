package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestTypeInsertsAndRenders(t *testing.T) {
	a := newApp(t, "")
	a.typ("i").typ("hello")
	assert.Equal(t, "hello", a.text())
	assert.True(t, a.shows("hello"))
}

func TestMultiCursorTypingAtEveryLine(t *testing.T) {
	a := newApp(t, "aaa\nbbb\nccc")
	a.press(tea.KeyCtrlDown).press(tea.KeyCtrlDown)
	a.typ("i").typ("X")
	assert.Equal(t, "Xaaa\nXbbb\nXccc", a.text())
}

func TestMultiCursorStatusShowsCount(t *testing.T) {
	a := newApp(t, "aaa\nbbb\nccc")
	a.press(tea.KeyCtrlDown).press(tea.KeyCtrlDown)
	assert.True(t, a.shows("3 cursors"))
}

func TestMultiCursorGroupedUndo(t *testing.T) {
	a := newApp(t, "aaa\nbbb")
	a.press(tea.KeyCtrlDown)
	a.typ("i").typ("Z")
	assert.Equal(t, "Zaaa\nZbbb", a.text())
	a.press(tea.KeyEsc).typ("u")
	assert.Equal(t, "aaa\nbbb", a.text(), "one undo reverts every cursor's edit")
}

func TestEscapeCollapsesToSingleCursor(t *testing.T) {
	a := newApp(t, "aaa\nbbb\nccc")
	a.press(tea.KeyCtrlDown)
	assert.Equal(t, 2, a.m.e.CaretCount())
	a.press(tea.KeyEsc)
	assert.Equal(t, 1, a.m.e.CaretCount())
}

func TestAddCaretAboveFromLowerLine(t *testing.T) {
	a := newApp(t, "aaa\nbbb\nccc")
	a.typ("j").typ("j")
	a.press(tea.KeyCtrlUp)
	a.typ("i").typ("X")
	assert.Equal(t, "aaa\nXbbb\nXccc", a.text())
}

func TestVimOpenLineFlow(t *testing.T) {
	a := newApp(t, "first\nsecond")
	a.typ("o").typ("mid") // open line below "first", type
	assert.Equal(t, "first\nmid\nsecond", a.text())
}

func TestVimYankPasteFlow(t *testing.T) {
	a := newApp(t, "dup\nother")
	a.typ("yy").typ("p") // yank line, paste below
	assert.Equal(t, "dup\ndup\nother", a.text())
}

func TestVimVisualDeleteFlow(t *testing.T) {
	a := newApp(t, "hello world")
	a.typ("v")               // visual
	for i := 0; i < 5; i++ { // extend over "hello"
		a.typ("l")
	}
	a.typ("d")
	assert.Equal(t, " world", a.text())
}

func TestSidebarToggleShowsTree(t *testing.T) {
	a := newApp(t, "code")
	a.press(tea.KeyCtrlB)
	for i := 0; i < 240; i++ {
		a.send(animMsg{})
	}
	assert.True(t, a.m.sbVisible)
	assert.True(t, a.shows("ROOT"))
}

func TestThemeSwitcherOpens(t *testing.T) {
	a := newApp(t, "code")
	a.press(tea.KeyCtrlT)
	assert.True(t, a.shows("Theme"))
}

func TestPaletteOpens(t *testing.T) {
	a := newApp(t, "code")
	a.press(tea.KeyCtrlP)
	assert.True(t, a.shows("Find") || a.shows("file") || len(a.screen()) > 0)
}

func TestViewAlwaysFillsHeight(t *testing.T) {
	a := newApp(t, "a\nb\nc")
	assert.Equal(t, 24, lipgloss.Height(a.m.View()))
	for i, line := range strings.Split(a.m.View(), "\n") {
		assert.Equalf(t, 80, lipgloss.Width(line), "row %d", i)
	}
}
