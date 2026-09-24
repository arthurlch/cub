package tui

import (
	"strings"

	"github.com/arthurlch/cub/internal/filetree"
	"github.com/arthurlch/cub/internal/icons"
	"github.com/arthurlch/cub/internal/input"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) toggleSidebar() {
	if m.tree == nil {
		return
	}
	if !m.sbVisible {
		m.sbVisible, m.sbFocused = true, true
		m.sbTarget, m.animating = 1, true
		return
	}
	m.sbFocused = !m.sbFocused
}

func (m *Model) hideSidebar() {
	m.sbVisible, m.sbFocused = false, false
	m.sbTarget, m.animating = 0, true
}

func (m *Model) sidebarKey(ev input.Event) {
	nodes := filetree.Flatten(m.tree)
	switch {
	case ev.Key == input.KeyUp, ev.Rune == 'k':
		if m.sbSelected > 0 {
			m.sbSelected--
		}
	case ev.Key == input.KeyDown, ev.Rune == 'j':
		if m.sbSelected < len(nodes)-1 {
			m.sbSelected++
		}
	case ev.Key == input.KeyEnter, ev.Key == input.KeyRight, ev.Rune == 'l':
		if len(nodes) > 0 {
			m.selectNode(nodes[m.sbSelected])
		}
	case ev.Key == input.KeyLeft, ev.Rune == 'h':
		if len(nodes) > 0 {
			if n := nodes[m.sbSelected]; n.IsDir && n.Expanded {
				n.Toggle()
			}
		}
	case ev.Rune == '>', ev.Rune == '+', ev.Rune == '=':
		m.sbWidth += 2
		m.clampSidebarWidth()
	case ev.Rune == '<', ev.Rune == '-':
		m.sbWidth -= 2
		m.clampSidebarWidth()
	case ev.Key == input.KeyEsc:
		m.hideSidebar()
	}

	count := len(filetree.Flatten(m.tree))
	if m.sbSelected >= count {
		m.sbSelected = count - 1
	}
	if m.sbSelected < 0 {
		m.sbSelected = 0
	}
}

func (m *Model) clampSidebarWidth() {
	const minWidth = 12
	max := m.width - 20
	if max < minWidth {
		max = minWidth
	}
	if m.sbWidth > max {
		m.sbWidth = max
	}
	if m.sbWidth < minWidth {
		m.sbWidth = minWidth
	}
}

func (m *Model) selectNode(n *filetree.Node) {
	if n.IsDir {
		n.Toggle()
		return
	}
	m.openInTab(n.Path)
	m.sbFocused = false
}

func (m Model) sidebarView(height int) string {
	w := m.sidebarDisplayWidth()
	if w <= 0 {
		return ""
	}
	if w < 3 {
		return lipgloss.NewStyle().Background(m.st.sbBg).Width(w).Height(height).Render("")
	}
	innerW := w - 1
	t := m.theme
	pad := lipgloss.NewStyle().Background(m.st.sbBg).Width(innerW)

	rows := make([]string, 0, height)
	rows = append(rows, pad.Render(""))
	rows = append(rows, m.st.sbHeader.Width(innerW).MaxWidth(innerW).Render("  "+strings.ToUpper(m.tree.Name)))
	rows = append(rows, pad.Render(""))

	nodes := filetree.Flatten(m.tree)
	capacity := height - 3
	if capacity < 0 {
		capacity = 0
	}
	top := 0
	if capacity > 0 && m.sbSelected >= capacity {
		top = m.sbSelected - capacity + 1
	}
	for i := 0; i < capacity; i++ {
		idx := top + i
		if idx >= len(nodes) {
			rows = append(rows, pad.Render(""))
			continue
		}
		n := nodes[idx]
		fg := t.SidebarFg
		if n.IsDir {
			fg = t.SidebarDir
		}
		rowBg := m.st.sbBg
		lead := " "
		leadFg := t.Accent
		if !n.IsDir && n.Path != "" && n.Path == m.e.FileName() {
			fg = t.Accent
		}
		if idx == m.sbSelected {
			rowBg = m.st.selBg
			if m.sbFocused {
				lead = "▎"
				fg = t.Fg
			} else {
				leadFg = m.st.dim
				fg = m.st.dim
			}
		}
		leadCell := lipgloss.NewStyle().Foreground(leadFg).Background(rowBg).Render(lead)
		body := lipgloss.NewStyle().Foreground(fg).Background(rowBg).Width(innerW - 1).MaxWidth(innerW - 1).Render(" " + truncate(nodeLabel(n), innerW-2))
		rows = append(rows, leadCell+body)
	}

	border := t.Accent
	if !m.sbFocused {
		border = m.st.dim
	}
	box := lipgloss.NewStyle().Background(m.st.sbBg).
		Border(lipgloss.RoundedBorder(), false, !m.sbRight, false, m.sbRight).
		BorderForeground(border).BorderBackground(m.theme.Bg)
	return box.Height(height).Render(strings.Join(rows, "\n"))
}

func nodeLabel(n *filetree.Node) string {
	indent := strings.Repeat("  ", n.Depth-1)
	if n.IsDir {
		arrow := "▸"
		if n.Expanded {
			arrow = "▾"
		}
		return indent + arrow + " " + icons.Folder(n.Expanded) + " " + n.Name
	}
	return indent + "  " + icons.File(n.Name) + " " + n.Name
}

func (m Model) helpView() string {
	lines := []string{
		"cub — keys",
		"",
		"i a         insert / append      o / O    open line below / above",
		"Esc         view mode            u / ^R   undo / redo",
		"h j k l     move                 w / b    word fwd / back",
		"0 ^ $       start/first/end      gg / G   top / bottom (<n>G line)",
		"v           visual select       y d x    yank / delete / cut",
		"dd yy       delete / yank line   p / P    paste after / before",
		"Ctrl+S      save                 Ctrl+Q   quit",
		"Ctrl+B      sidebar              Ctrl+P   fuzzy file finder",
		"Ctrl+T      switch theme         Ctrl+G   toggle terminal",
		"Ctrl+←/→    prev / next tab      Ctrl+W   close tab",
		"Ctrl+N      new tab              Ctrl+E   sidebar left / right",
		"Ctrl+↑/↓    add cursor up/down   Ctrl+H   this help",
		"",
		"terminal    Ctrl+G toggle · Ctrl+N new · Ctrl+←/→ switch · Esc editor · Ctrl+↑/↓ resize",
		"sidebar     Ctrl+B focus/return   Esc close   j/k move   Enter open   < > resize",
		"mouse       click move · drag select · click ✕ close tab · drag borders resize",
		"",
		"config      ~/.config/cub/config.json — theme, sidebar side, keybindings",
		"",
		"press any key to close",
	}
	return m.st.modal.Render(strings.Join(lines, "\n"))
}

func truncate(s string, w int) string {
	r := []rune(s)
	if len(r) > w {
		return string(r[:w])
	}
	return s
}
