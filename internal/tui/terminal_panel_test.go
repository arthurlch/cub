package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestKeyToBytes(t *testing.T) {
	cases := []struct {
		msg  tea.KeyMsg
		want []byte
	}{
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("ab")}, []byte("ab")},
		{tea.KeyMsg{Type: tea.KeyEnter}, []byte{'\r'}},
		{tea.KeyMsg{Type: tea.KeyBackspace}, []byte{0x7f}},
		{tea.KeyMsg{Type: tea.KeyEsc}, []byte{0x1b}},
		{tea.KeyMsg{Type: tea.KeyUp}, []byte("\x1b[A")},
		{tea.KeyMsg{Type: tea.KeyCtrlC}, []byte{3}},
		{tea.KeyMsg{Type: tea.KeyCtrlD}, []byte{4}},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, keyToBytes(c.msg))
	}
}

func TestTerminalShrinksEditorViewport(t *testing.T) {
	a := newApp(t, "code")
	full := a.m.termHeight()
	assert.Equal(t, 0, full, "terminal starts hidden")

	a.m.termVisible = true
	a.m.applySize()
	assert.Greater(t, a.m.termHeight(), 0)
	assert.Equal(t, 80, lipgloss.Width(a.m.View()))
}

func TestTerminalResizeByKeyAndDrag(t *testing.T) {
	a := newApp(t, "code")
	a.m.termVisible = true
	a.m.applySize()
	base := a.m.termHeight()

	a.m.resizeTerminal(3)
	assert.Equal(t, base+3, a.m.termHeight(), "grow adds rows")

	a.m.resizeTerminal(-100)
	assert.GreaterOrEqual(t, a.m.termHeight(), 4, "shrink clamps to a minimum")

	a.m.setTerminalTop(15) // panel spans from y=15 down; height 24
	assert.Equal(t, 24-1-15, a.m.termHeight(), "dragging the header sets the panel top")
}

func TestTerminalKeyGrowsWhenFocused(t *testing.T) {
	if testing.Short() {
		t.Skip("skips spawning a shell in short mode")
	}
	a := newApp(t, "code")
	a.m.toggleTerminal()
	if a.m.activeTerm() == nil {
		t.Skip("no pty available in this environment")
	}
	defer a.m.closeAllTerminals()

	before := a.m.termHeight()
	a.send(tea.KeyMsg{Type: tea.KeyCtrlUp})
	assert.Greater(t, a.m.termHeight(), before, "ctrl+up grows the focused terminal")
}

func TestMultipleTerminalsAndSwitch(t *testing.T) {
	if testing.Short() {
		t.Skip("skips spawning shells in short mode")
	}
	a := newApp(t, "code")
	a.m.toggleTerminal()
	if a.m.activeTerm() == nil {
		t.Skip("no pty available in this environment")
	}
	defer a.m.closeAllTerminals()

	a.m.spawnTerminal()
	a.m.spawnTerminal()
	assert.Equal(t, 3, len(a.m.terms), "spawned three terminals")
	assert.Equal(t, 2, a.m.termActive, "newest terminal is active")

	a.m.switchTerminal(1)
	assert.Equal(t, 0, a.m.termActive, "next wraps around")
	a.m.switchTerminal(-1)
	assert.Equal(t, 2, a.m.termActive, "prev wraps around")

	a.m.focusTerminal(1)
	assert.Equal(t, 1, a.m.termActive, "focusing a terminal tab selects it")
}

func TestCloseTerminals(t *testing.T) {
	if testing.Short() {
		t.Skip("skips spawning shells in short mode")
	}
	a := newApp(t, "code")
	a.m.toggleTerminal()
	if a.m.activeTerm() == nil {
		t.Skip("no pty available in this environment")
	}
	defer a.m.closeAllTerminals()

	a.m.spawnTerminal()
	a.m.spawnTerminal()
	a.m.focusTerminal(1)
	assert.Equal(t, 3, len(a.m.terms))

	a.m.closeTerminalAt(0)
	assert.Equal(t, 2, len(a.m.terms), "closing a terminal removes it")
	assert.Equal(t, 0, a.m.termActive, "active index follows when a lower tab closes")

	a.m.dropActiveTerminal()
	a.m.dropActiveTerminal()
	assert.Equal(t, 0, len(a.m.terms), "closing the last terminal empties the panel")
	assert.False(t, a.m.termVisible, "no terminals left hides the panel")
}

func TestHiddenTerminalDoesNotCaptureKeys(t *testing.T) {
	a := newApp(t, "abc")
	a.m.termVisible = false
	a.m.termFocused = true // stale focus flag must not route keys to a hidden panel
	a.typ("i").typ("Z")
	assert.Equal(t, "Zabc", a.text(), "keys reach the editor when the terminal is hidden")
}

func TestTerminalGenerationGuardsStaleUpdates(t *testing.T) {
	a := newApp(t, "code")
	a.m.termGen = 5
	tm, cmd := a.m.Update(termMsg{gen: 2})
	a.m = tm.(Model)
	assert.Nil(t, cmd, "stale terminal update (old gen) is dropped")
}

func TestTerminalToggleSpawnsAndRenders(t *testing.T) {
	if testing.Short() {
		t.Skip("skips spawning a shell in short mode")
	}
	a := newApp(t, "code")
	cmd := a.m.toggleTerminal()
	if a.m.activeTerm() == nil {
		t.Skip("no pty available in this environment")
	}
	defer a.m.closeAllTerminals()

	assert.True(t, a.m.termVisible)
	assert.True(t, a.m.termFocused)
	assert.NotNil(t, cmd)
	assert.True(t, a.shows("TERMINAL"))

	a.m.toggleTerminal()
	assert.False(t, a.m.termVisible)
}
