package tui

import (
	"github.com/arthurlch/cub/internal/filetree"
	tea "github.com/charmbracelet/bubbletea"
)

// NOTE: i don't plan much more mouse interaction than this,
// if we want to do more mouse stuff, we should probably
// refactor the model to have a more general "click" and "drag" handling system.
// But it the idea of using the mouse is just not the point of this editor,
// so i don't think it's worth the effort.
// The mouse handling is mostly just to make it easier to use the sidebar and terminal.
// for the new users

func (m *Model) handleMouse(msg tea.MouseMsg) tea.Cmd {
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		m.wheel(-1, msg.X)
		return nil
	case tea.MouseButtonWheelDown:
		m.wheel(1, msg.X)
		return nil
	}

	switch msg.Action {
	case tea.MouseActionPress:
		if msg.Button != tea.MouseButtonLeft {
			return nil
		}
		if msg.Y == 0 {
			m.clickTab(msg.X)
			return nil
		}
		if m.termVisible {
			hr := m.termHeaderRow()
			if msg.Y == hr {
				return m.clickTerminalHeader(msg.X)
			}
			if msg.Y > hr && msg.Y < m.height-1 {
				return m.focusTerminal(m.termActive)
			}
		}
		sb := m.sidebarDisplayWidth()
		if sb > 0 {
			border := m.sidebarBorderX()
			if msg.X == border {
				m.resizing = true
				return nil
			}
			if m.inSidebarX(msg.X) {
				m.sbFocused = true
				m.clickSidebar(msg.Y)
				return nil
			}
		}
		m.sbFocused = false
		m.termFocused = false
		m.e.EndSelection()
		m.moveCursorTo(msg.X, msg.Y)
		m.dragging = true

	case tea.MouseActionMotion:
		if m.termResizing {
			m.setTerminalTop(msg.Y)
			return nil
		}
		if m.resizing {
			if m.sbRight {
				m.sbWidth = m.width - msg.X
			} else {
				m.sbWidth = msg.X + 1
			}
			m.clampSidebarWidth()
			m.applySize()
			return nil
		}
		if m.dragging {
			if active, _, _, _, _ := m.e.Selection(); !active {
				m.e.StartSelection()
			}
			m.moveCursorTo(msg.X, msg.Y)
		}

	case tea.MouseActionRelease:
		m.resizing = false
		m.dragging = false
		m.termResizing = false
	}
	return nil
}

func (m *Model) clickTerminalHeader(x int) tea.Cmd {
	for _, tab := range m.termTabLayout() {
		if x >= tab.x && x < tab.x+tab.w {
			if tab.add {
				return m.spawnTerminal()
			}
			if x >= tab.x+tab.w-2 {
				return m.closeTerminalAt(tab.idx)
			}
			return m.focusTerminal(tab.idx)
		}
	}
	m.termResizing = true
	return nil
}

func (m *Model) focusTerminal(i int) tea.Cmd {
	m.termFocused = true
	if i < 0 || i >= len(m.terms) || i == m.termActive {
		return nil
	}
	m.termActive = i
	m.termGen++
	if s := m.activeTerm(); s != nil {
		s.Resize(m.termCols(), m.termRows())
	}
	return m.waitTerm()
}

func (m *Model) setTerminalTop(y int) {
	m.termLines = m.height - 1 - y
	m.applySize()
	if s := m.activeTerm(); s != nil && m.termVisible {
		s.Resize(m.termCols(), m.termRows())
	}
}

func (m Model) sidebarBorderX() int {
	if m.sbRight {
		return m.width - m.sidebarDisplayWidth()
	}
	return m.sidebarDisplayWidth() - 1
}

func (m Model) inSidebarX(x int) bool {
	sb := m.sidebarDisplayWidth()
	if sb <= 0 {
		return false
	}
	if m.sbRight {
		return x > m.width-sb
	}
	return x < sb-1
}

func (m *Model) moveCursorTo(x, y int) {
	textX := x - m.editorOriginX() - gutterWidth
	if textX < 0 {
		textX = 0
	}
	row := y - 1
	if row < 0 {
		row = 0
	}
	offRow, offCol := m.e.Offset()
	m.e.ClearCarets()
	m.e.MoveTo(offRow+row, offCol+textX)
	m.applySize()
}

func (m *Model) clickSidebar(y int) {
	nodes := filetree.Flatten(m.tree)
	capacity := m.height - 5
	if capacity < 1 {
		return
	}
	top := 0
	if m.sbSelected >= capacity {
		top = m.sbSelected - capacity + 1
	}
	j := y - 4
	idx := top + j
	if j >= 0 && j < capacity && idx >= 0 && idx < len(nodes) {
		m.sbSelected = idx
		m.selectNode(nodes[idx])
	}
}

func (m *Model) wheel(dir, x int) {
	if m.inSidebarX(x) || (m.sidebarDisplayWidth() > 0 && x == m.sidebarBorderX()) {
		nodes := filetree.Flatten(m.tree)
		m.sbSelected += dir
		if m.sbSelected < 0 {
			m.sbSelected = 0
		}
		if m.sbSelected >= len(nodes) {
			m.sbSelected = len(nodes) - 1
		}
		return
	}
	for i := 0; i < 3; i++ {
		if dir < 0 {
			m.e.MoveUp()
		} else {
			m.e.MoveDown()
		}
	}
	m.applySize()
}
