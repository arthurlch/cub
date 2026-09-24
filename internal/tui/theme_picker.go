package tui

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/arthurlch/cub/internal/input"
	"github.com/arthurlch/cub/internal/theme"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) themeKey(ev input.Event) {
	rune := !ev.Ctrl && !ev.Alt
	switch {
	case ev.Key == input.KeyUp, rune && ev.Rune == 'k':
		if m.themeSel > 0 {
			m.themeSel--
		}
		m.applyTheme(m.themeSel)
	case ev.Key == input.KeyDown, rune && ev.Rune == 'j':
		if m.themeSel < len(theme.Themes)-1 {
			m.themeSel++
		}
		m.applyTheme(m.themeSel)
	case ev.Key == input.KeyEnter, rune && ev.Rune == 'l':
		m.themeIdx = m.themeSel
		m.showTheme = false
		m.saveConfig()
	case ev.Key == input.KeyEsc:
		m.applyTheme(m.themeIdx)
		m.showTheme = false
	}
}

func (m *Model) applyTheme(i int) {
	m.theme = theme.Themes[i]
	m.st = buildStyles(m.theme)
}

func (m Model) themePickerView() string {
	t := m.theme
	sel := lipgloss.NewStyle().Foreground(t.CursorFg).Background(t.Accent).Bold(true)
	item := lipgloss.NewStyle().Foreground(t.Fg).Background(t.StatusBg)
	swatch := func(th theme.Theme) string {
		s := lipgloss.NewStyle().Background(th.Bg)
		dot := func(clr lipgloss.Color) string {
			return lipgloss.NewStyle().Foreground(clr).Background(th.Bg).Render("●")
		}
		return s.Render(" ") +
			dot(th.Token(chroma.Keyword)) +
			dot(th.Token(chroma.LiteralString)) +
			dot(th.Token(chroma.NameFunction)) +
			dot(th.Token(chroma.LiteralNumber)) +
			s.Render(" ")
	}

	const visible = 14
	start := 0
	if m.themeSel >= visible {
		start = m.themeSel - visible + 1
	}
	end := start + visible
	if end > len(theme.Themes) {
		end = len(theme.Themes)
	}

	header := fmt.Sprintf("Theme  %d/%d", m.themeSel+1, len(theme.Themes))
	rows := []string{m.st.sbHeader.Background(t.StatusBg).Render(header), ""}
	for i := start; i < end; i++ {
		th := theme.Themes[i]
		row := "  " + th.Name
		if i == m.themeSel {
			row = "▸ " + th.Name
		}
		for lipgloss.Width(row) < 22 {
			row += " "
		}
		line := row + swatch(th)
		style := item
		if i == m.themeSel {
			style = sel
		}
		rows = append(rows, style.Render(line))
	}
	rows = append(rows, "", m.st.hint.Render("↑/↓ preview · Enter apply · Esc cancel"))

	return m.st.modal.Render(strings.Join(rows, "\n"))
}
