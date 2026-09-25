package tui

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/arthurlch/cub/internal/icons"
	"github.com/arthurlch/cub/internal/input"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

const paletteRows = 12

var errStopWalk = errors.New("stop")

func (m *Model) openPalette() {
	m.showPalette = true
	m.paletteQuery = ""
	m.paletteSel = 0
	if m.paletteFiles == nil {
		m.paletteFiles = collectFiles(".")
	}
	m.recomputeMatches()
}

func collectFiles(root string) []string {
	skip := map[string]bool{".git": true, "node_modules": true, "vendor": true, ".idea": true, "dist": true}
	var out []string
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if path != root && skip[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if rel, e := filepath.Rel(root, path); e == nil {
			out = append(out, rel)
		}
		if len(out) >= 5000 {
			return errStopWalk
		}
		return nil
	})
	return out
}

func (m *Model) recomputeMatches() {
	m.paletteMatches = m.paletteMatches[:0]
	if m.paletteQuery == "" {
		for i := range m.paletteFiles {
			if len(m.paletteMatches) >= 200 {
				break
			}
			m.paletteMatches = append(m.paletteMatches, i)
		}
	} else {
		for _, r := range fuzzy.Find(m.paletteQuery, m.paletteFiles) {
			m.paletteMatches = append(m.paletteMatches, r.Index)
		}
	}
	if m.paletteSel >= len(m.paletteMatches) {
		m.paletteSel = len(m.paletteMatches) - 1
	}
	if m.paletteSel < 0 {
		m.paletteSel = 0
	}
}

func (m *Model) paletteKey(ev input.Event) {
	switch ev.Key {
	case input.KeyEsc:
		m.showPalette = false
	case input.KeyUp:
		if m.paletteSel > 0 {
			m.paletteSel--
		}
	case input.KeyDown:
		if m.paletteSel < len(m.paletteMatches)-1 {
			m.paletteSel++
		}
	case input.KeyEnter:
		if m.paletteSel >= 0 && m.paletteSel < len(m.paletteMatches) {
			m.openInTab(m.paletteFiles[m.paletteMatches[m.paletteSel]])
		}
		m.showPalette = false
	case input.KeyBackspace:
		if r := []rune(m.paletteQuery); len(r) > 0 {
			m.paletteQuery = string(r[:len(r)-1])
			m.recomputeMatches()
		}
	default:
		if ev.Rune != 0 && !ev.Ctrl && !ev.Alt {
			m.paletteQuery += string(ev.Rune)
			m.recomputeMatches()
		}
	}
}

func (m Model) paletteView() string {
	t := m.theme
	width := 52
	prompt := lipgloss.NewStyle().Foreground(t.Accent).Background(t.StatusBg).Bold(true)
	field := lipgloss.NewStyle().Foreground(t.Fg).Background(t.StatusBg)
	rule := lipgloss.NewStyle().Foreground(t.Border).Background(t.StatusBg)
	selStyle := lipgloss.NewStyle().Foreground(t.CursorFg).Background(t.Accent).Bold(true)
	item := lipgloss.NewStyle().Foreground(t.Fg).Background(t.StatusBg)
	dim := lipgloss.NewStyle().Foreground(m.st.dim).Background(t.StatusBg)

	rows := []string{
		prompt.Render("  ") + field.Width(width-2).Render(m.paletteQuery+"▌"),
		rule.Render(strings.Repeat("─", width)),
	}

	start := 0
	if m.paletteSel >= paletteRows {
		start = m.paletteSel - paletteRows + 1
	}
	shown := 0
	for i := start; i < len(m.paletteMatches) && shown < paletteRows; i++ {
		rel := m.paletteFiles[m.paletteMatches[i]]
		label := " " + icons.File(rel) + " " + rel
		style := item
		if i == m.paletteSel {
			style = selStyle
		}
		rows = append(rows, style.Width(width).MaxWidth(width).Render(truncate(label, width)))
		shown++
	}
	if len(m.paletteMatches) == 0 {
		rows = append(rows, dim.Render(" no matches"))
	}
	rows = append(rows, "", m.st.hint.Render("↑/↓ select · Enter open · Esc cancel"))

	return m.st.modal.Render(strings.Join(rows, "\n"))
}
