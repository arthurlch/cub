package tui

import (
	"fmt"
	"strings"

	"github.com/arthurlch/cub/internal/input"
	"github.com/arthurlch/cub/internal/keymap"
	"github.com/arthurlch/cub/internal/terminal"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func (m *Model) activeTerm() *terminal.Session {
	if m.termActive < 0 || m.termActive >= len(m.terms) {
		return nil
	}
	return m.terms[m.termActive]
}

func (m Model) termHeight() int {
	if !m.termVisible {
		return 0
	}
	h := m.termLines
	if h == 0 {
		h = m.height / 3
	}
	min := 4
	max := m.height - 6
	if max < min {
		max = min
	}
	if h < min {
		h = min
	}
	if h > max {
		h = max
	}
	return h
}

func (m Model) termHeaderRow() int {
	return m.height - 1 - m.termHeight()
}

func (m Model) termCols() int {
	if m.width < 1 {
		return 1
	}
	return m.width
}

func (m Model) termRows() int {
	if r := m.termHeight() - 1; r > 0 {
		return r
	}
	return 1
}

func (m *Model) resizeTerminal(delta int) {
	if !m.termVisible {
		return
	}
	m.termLines = m.termHeight() + delta
	m.applySize()
	if s := m.activeTerm(); s != nil {
		s.Resize(m.termCols(), m.termRows())
	}
}

func (m *Model) toggleTerminal() tea.Cmd {
	if m.termVisible {
		if m.termFocused {
			m.termVisible, m.termFocused = false, false
			m.applySize()
			return nil
		}
		m.termFocused = true
		return nil
	}

	m.termVisible, m.termFocused = true, true
	m.applySize()

	if s := m.activeTerm(); s != nil && !s.Closed() {
		s.Resize(m.termCols(), m.termRows())
		m.termGen++
		return m.waitTerm()
	}
	return m.spawnTerminal()
}

func (m *Model) spawnTerminal() tea.Cmd {
	s, err := terminal.New(m.termCols(), m.termRows())
	if err != nil {
		m.e.ShowMessage("terminal unavailable: " + err.Error())
		if len(m.terms) == 0 {
			m.termVisible, m.termFocused = false, false
			m.applySize()
		}
		return nil
	}
	m.terms = append(m.terms, s)
	m.termActive = len(m.terms) - 1
	m.termGen++
	return m.waitTerm()
}

func (m *Model) switchTerminal(delta int) tea.Cmd {
	n := len(m.terms)
	if n < 2 {
		return nil
	}
	m.termActive = (m.termActive + delta + n) % n
	m.termGen++
	if s := m.activeTerm(); s != nil {
		s.Resize(m.termCols(), m.termRows())
	}
	return m.waitTerm()
}

func (m *Model) dropActiveTerminal() tea.Cmd {
	return m.closeTerminalAt(m.termActive)
}

func (m *Model) closeTerminalAt(i int) tea.Cmd {
	if i < 0 || i >= len(m.terms) {
		return nil
	}
	m.terms[i].Close()
	m.terms = append(m.terms[:i], m.terms[i+1:]...)
	if len(m.terms) == 0 {
		m.termActive = 0
		m.termVisible, m.termFocused = false, false
		m.applySize()
		return nil
	}
	if i < m.termActive {
		m.termActive--
	}
	if m.termActive >= len(m.terms) {
		m.termActive = len(m.terms) - 1
	}
	m.termGen++
	if s := m.activeTerm(); s != nil {
		s.Resize(m.termCols(), m.termRows())
	}
	return m.waitTerm()
}

func (m *Model) closeAllTerminals() {
	for _, s := range m.terms {
		s.Close()
	}
	m.terms = nil
}

func (m Model) waitTerm() tea.Cmd {
	s := m.activeTerm()
	gen := m.termGen
	if s == nil {
		return nil
	}
	return func() tea.Msg {
		<-s.Updates()
		return termMsg{gen: gen}
	}
}

func (m *Model) terminalKey(msg tea.KeyMsg, ev input.Event) (tea.Model, tea.Cmd) {
	if ev.Key == input.KeyEsc {
		m.termFocused = false
		return *m, nil
	}
	switch m.km.Global(ev) {
	case keymap.Quit:
		m.closeAllTerminals()
		return *m, tea.Quit
	case keymap.Terminal:
		cmd := m.toggleTerminal()
		return *m, cmd
	case keymap.CursorAddAbove:
		m.resizeTerminal(1)
		return *m, nil
	case keymap.CursorAddBelow:
		m.resizeTerminal(-1)
		return *m, nil
	case keymap.TabNew:
		cmd := m.spawnTerminal()
		return *m, cmd
	case keymap.TabClose:
		cmd := m.dropActiveTerminal()
		return *m, cmd
	case keymap.TabNext:
		cmd := m.switchTerminal(1)
		return *m, cmd
	case keymap.TabPrev:
		cmd := m.switchTerminal(-1)
		return *m, cmd
	}
	if b := keyToBytes(msg); b != nil {
		m.activeTerm().Write(b)
	}
	return *m, nil
}

type termTab struct {
	x, w, idx int
	add       bool
}

func (m Model) termTabLayout() []termTab {
	var tabs []termTab
	x := lipgloss.Width(m.termLabel())
	for i := range m.terms {
		w := lipgloss.Width(m.termTabLabel(i))
		tabs = append(tabs, termTab{x: x, w: w, idx: i})
		x += w
	}
	tabs = append(tabs, termTab{x: x, w: 3, idx: -1, add: true})
	return tabs
}

func (m Model) termLabel() string         { return " ▚ TERMINAL " }
func (m Model) termTabLabel(i int) string { return fmt.Sprintf(" %d ✕ ", i+1) }

func (m Model) terminalView(width, height int) string {
	header := m.terminalHeader(width)
	rows := height - 1
	if rows < 0 {
		rows = 0
	}
	body := ""
	if s := m.activeTerm(); s != nil {
		body = s.Render()
	}
	return lipgloss.JoinVertical(lipgloss.Left, header, fitPanel(body, width, rows))
}

func (m Model) terminalHeader(width int) string {
	bg := m.st.sbBg
	accent := m.theme.Accent
	if !m.termFocused {
		accent = m.st.dim
	}
	labelStyle := lipgloss.NewStyle().Foreground(m.theme.CursorFg).Background(accent).Bold(true)
	barStyle := lipgloss.NewStyle().Foreground(m.theme.Fg).Background(bg)
	tabStyle := lipgloss.NewStyle().Foreground(m.st.dim).Background(bg)
	activeTab := lipgloss.NewStyle().Foreground(m.theme.Fg).Background(m.st.selBg).Bold(true)
	hintStyle := lipgloss.NewStyle().Foreground(m.st.dim).Background(bg)

	var b strings.Builder
	b.WriteString(labelStyle.Render(m.termLabel()))
	for i := range m.terms {
		label := m.termTabLabel(i)
		if i == m.termActive {
			b.WriteString(activeTab.Render(label))
		} else {
			b.WriteString(tabStyle.Render(label))
		}
	}
	b.WriteString(tabStyle.Render(" + "))

	hint := " Ctrl+N new · Ctrl+W close · Ctrl+←/→ switch · Esc editor · Ctrl+G hide "
	if !m.termFocused {
		hint = " Ctrl+G / click to focus "
	}
	used := lipgloss.Width(b.String()) + lipgloss.Width(hint)
	fill := width - used
	if fill < 0 {
		fill = 0
		hint = ""
	}
	b.WriteString(barStyle.Render(strings.Repeat(" ", fill)))
	b.WriteString(hintStyle.Render(hint))
	return fitPanel(b.String(), width, 1)
}

func fitPanel(body string, width, rows int) string {
	src := strings.Split(body, "\n")
	out := make([]string, rows)
	for i := 0; i < rows; i++ {
		line := ""
		if i < len(src) {
			line = src[i]
		}
		w := lipgloss.Width(line)
		switch {
		case w > width:
			line = ansi.Truncate(line, width, "")
		case w < width:
			line += strings.Repeat(" ", width-w)
		}
		out[i] = line
	}
	return strings.Join(out, "\n")
}

func keyToBytes(msg tea.KeyMsg) []byte {
	switch msg.Type {
	case tea.KeyRunes:
		return []byte(string(msg.Runes))
	case tea.KeySpace:
		return []byte{' '}
	case tea.KeyEnter:
		return []byte{'\r'}
	case tea.KeyTab:
		return []byte{'\t'}
	case tea.KeyBackspace:
		return []byte{0x7f}
	case tea.KeyEsc:
		return []byte{0x1b}
	case tea.KeyUp:
		return []byte("\x1b[A")
	case tea.KeyDown:
		return []byte("\x1b[B")
	case tea.KeyRight:
		return []byte("\x1b[C")
	case tea.KeyLeft:
		return []byte("\x1b[D")
	case tea.KeyHome:
		return []byte("\x1b[H")
	case tea.KeyEnd:
		return []byte("\x1b[F")
	case tea.KeyPgUp:
		return []byte("\x1b[5~")
	case tea.KeyPgDown:
		return []byte("\x1b[6~")
	case tea.KeyDelete:
		return []byte("\x1b[3~")
	}
	if msg.Type >= tea.KeyCtrlA && msg.Type <= tea.KeyCtrlZ {
		return []byte{byte(msg.Type-tea.KeyCtrlA) + 1}
	}
	return nil
}
