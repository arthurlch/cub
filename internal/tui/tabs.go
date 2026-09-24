package tui

import (
	"github.com/arthurlch/cub/internal/editor"
	"github.com/arthurlch/cub/internal/icons"
	"github.com/charmbracelet/lipgloss"
)

// tabs kinda hacky now but hey it works !
// NOTE: Will need to be refactored when we add more features like tab
// reordering, closing tabs with the mouse, etc.
// okay I can close tabwith with mouse awesome ..

type tabBox struct {
	x, w, doc int
}

func (m *Model) setActive(i int) {
	if i < 0 {
		i = 0
	}
	if i >= len(m.docs) {
		i = len(m.docs) - 1
	}
	m.active = i
	m.e = m.docs[i]
	m.closeConfirm = false
	m.applySize()
}

func (m *Model) isScratch(e *editor.Editor) bool {
	return e.FileName() == "" && !e.Modified() && e.LineCount() == 1 && len(e.Line(0)) == 0
}

func (m *Model) openInTab(path string) {
	for i, d := range m.docs {
		if d.FileName() == path {
			m.setActive(i)
			return
		}
	}
	if m.isScratch(m.e) {
		if err := m.e.Open(path); err != nil {
			m.e.ShowMessage("open failed: " + err.Error())
		}
		m.applySize()
		return
	}
	nd := editor.New()
	if err := nd.Open(path); err != nil {
		nd.ShowMessage("open failed: " + err.Error())
	}
	m.docs = append(m.docs, nd)
	m.setActive(len(m.docs) - 1)
}

func (m *Model) requestClose() {
	if m.e.Modified() && !m.closeConfirm {
		m.closeConfirm = true
		m.e.ShowMessage("unsaved changes — Ctrl+W again to close")
		return
	}
	m.closeTab()
}

func (m *Model) closeTab() {
	if len(m.docs) <= 1 {
		m.docs = []*editor.Editor{editor.New()}
		m.setActive(0)
		return
	}
	i := m.active
	m.docs = append(m.docs[:i], m.docs[i+1:]...)
	m.setActive(i)
}

func (m *Model) newTab() {
	m.docs = append(m.docs, editor.New())
	m.setActive(len(m.docs) - 1)
}

func (m *Model) nextTab() { m.setActive((m.active + 1) % len(m.docs)) }
func (m *Model) prevTab() { m.setActive((m.active - 1 + len(m.docs)) % len(m.docs)) }

func (m Model) tabLabel(d *editor.Editor) string {
	name := baseName(d.FileName())
	if d.Modified() {
		name += " ●"
	}
	return " " + icons.File(d.FileName()) + " " + name + " ✕ "
}

func (m Model) tabLayout(width int) []tabBox {
	lw := func(d *editor.Editor) int { return lipgloss.Width(m.tabLabel(d)) }

	total := 0
	for _, d := range m.docs {
		total += lw(d)
	}

	start := 0
	if total > width {
		w := 0
		for i := m.active; i >= 0; i-- {
			w += lw(m.docs[i])
			if w > width {
				start = i + 1
				break
			}
		}
	}

	var boxes []tabBox
	x := 0
	for i := start; i < len(m.docs); i++ {
		w := lw(m.docs[i])
		if x+w > width {
			break
		}
		boxes = append(boxes, tabBox{x: x, w: w, doc: i})
		x += w
	}
	return boxes
}

func (m *Model) clickTab(x int) {
	for _, bx := range m.tabLayout(m.width) {
		if x >= bx.x && x < bx.x+bx.w {
			m.setActive(bx.doc)
			if x >= bx.x+bx.w-2 {
				m.requestClose()
			}
			return
		}
	}
}
