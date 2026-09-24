package tui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/alecthomas/chroma/v2"
	"github.com/arthurlch/cub/internal/config"
	"github.com/arthurlch/cub/internal/document"
	"github.com/arthurlch/cub/internal/editor"
	"github.com/arthurlch/cub/internal/filetree"
	"github.com/arthurlch/cub/internal/icons"
	"github.com/arthurlch/cub/internal/keymap"
	"github.com/arthurlch/cub/internal/syntax"
	"github.com/arthurlch/cub/internal/terminal"
	"github.com/arthurlch/cub/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// charmbracelet is carrying my a$$ for the TUI and it works well actually

type animMsg struct{}

type termMsg struct{ gen int }

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg { return animMsg{} })
}

const (
	gutterWidth = 5
	tabStop     = 4
)

type styleSet struct {
	sbBg  lipgloss.Color
	selBg lipgloss.Color
	dim   lipgloss.Color

	gutter    lipgloss.Style
	gutterHi  lipgloss.Style
	cursor    lipgloss.Style
	statusBar lipgloss.Style
	modePill  lipgloss.Style
	plArrow   lipgloss.Style
	title     lipgloss.Style
	sbHeader  lipgloss.Style
	modal     lipgloss.Style
	hint      lipgloss.Style
}

func buildStyles(t theme.Theme) styleSet {
	sbBg := lerpColor(t.Bg, t.Fg, 0.1)
	selBg := lerpColor(t.Bg, t.Accent, 0.28)
	dim := lerpColor(t.Fg, t.Bg, 0.45)
	return styleSet{
		sbBg: sbBg, selBg: selBg, dim: dim,
		gutter:    lipgloss.NewStyle().Foreground(t.Gutter).Background(t.Bg),
		gutterHi:  lipgloss.NewStyle().Foreground(t.GutterHi).Background(t.CurrentLine),
		cursor:    lipgloss.NewStyle().Foreground(t.CursorFg).Background(t.CursorBg),
		statusBar: lipgloss.NewStyle().Foreground(t.StatusFg).Background(t.StatusBg),
		modePill:  lipgloss.NewStyle().Foreground(t.CursorFg).Background(t.Accent).Bold(true),
		plArrow:   lipgloss.NewStyle().Foreground(t.Accent).Background(t.StatusBg),
		title:     lipgloss.NewStyle().Background(t.Bg),
		sbHeader:  lipgloss.NewStyle().Foreground(dim).Background(sbBg).Bold(true),
		modal:     lipgloss.NewStyle().Foreground(t.Fg).Background(t.StatusBg).Border(lipgloss.RoundedBorder()).BorderForeground(t.Accent).Padding(1, 2),
		hint:      lipgloss.NewStyle().Foreground(dim).Background(t.StatusBg),
	}
}

type Model struct {
	e      *editor.Editor
	docs   []*editor.Editor
	active int
	km     *keymap.Keymap

	cfg     config.Config
	persist bool

	closeConfirm bool

	tree       *filetree.Node
	sbVisible  bool
	sbFocused  bool
	sbRight    bool
	sbSelected int
	sbWidth    int

	theme    theme.Theme
	themeIdx int
	st       styleSet

	showHelp  bool
	showTheme bool
	themeSel  int

	showPalette    bool
	paletteQuery   string
	paletteFiles   []string
	paletteMatches []int
	paletteSel     int

	terms        []*terminal.Session
	termActive   int
	termGen      int
	termVisible  bool
	termFocused  bool
	termLines    int
	termResizing bool

	resizing bool
	dragging bool
	width    int
	height   int

	spring    harmonica.Spring
	sbReveal  float64
	sbVel     float64
	sbTarget  float64
	animating bool
}

func New(e *editor.Editor, tree *filetree.Node, sidebarWidth int) Model {
	t := theme.Default()
	return Model{
		e: e, docs: []*editor.Editor{e}, km: keymap.Default(), tree: tree, sbWidth: sidebarWidth,
		theme: t, st: buildStyles(t),
		spring: harmonica.NewSpring(harmonica.FPS(60), 14.0, 1.0),
	}
}

func (m Model) WithConfig(cfg config.Config, km *keymap.Keymap) Model {
	m.cfg = cfg
	m.km = km
	m.persist = true
	m.sbRight = cfg.SidebarRight
	for i, th := range theme.Themes {
		if th.Name == cfg.Theme {
			m.themeIdx, m.theme, m.st = i, th, buildStyles(th)
			break
		}
	}
	return m
}

func (m *Model) saveConfig() {
	if !m.persist {
		return
	}
	m.cfg.Theme = theme.Themes[m.themeIdx].Name
	m.cfg.SidebarRight = m.sbRight
	_ = config.Save(m.cfg)
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.applySize()
		if s := m.activeTerm(); s != nil && m.termVisible {
			s.Resize(m.termCols(), m.termRows())
		}
	case animMsg:
		m.sbReveal, m.sbVel = m.spring.Update(m.sbReveal, m.sbVel, m.sbTarget)
		if math.Abs(m.sbReveal-m.sbTarget) < 0.002 && math.Abs(m.sbVel) < 0.002 {
			m.sbReveal, m.sbVel, m.animating = m.sbTarget, 0, false
		}
		m.applySize()
		if m.animating {
			return m, tickCmd()
		}
	case termMsg:
		if msg.gen != m.termGen {
			return m, nil
		}
		s := m.activeTerm()
		if s == nil {
			return m, nil
		}
		if s.Closed() {
			return m, m.dropActiveTerminal()
		}
		return m, m.waitTerm()
	case tea.MouseMsg:
		if !m.showTheme && !m.showHelp && !m.showPalette {
			return m, m.handleMouse(msg)
		}
	case tea.KeyMsg:
		ev := fromTea(msg)
		if m.showTheme {
			m.themeKey(ev)
			return m, nil
		}
		if m.showPalette {
			m.paletteKey(ev)
			return m, nil
		}
		if m.termVisible && m.termFocused && m.activeTerm() != nil {
			return m.terminalKey(msg, ev)
		}
		if action := m.km.Global(ev); action != "" {
			if action == keymap.Quit {
				m.closeAllTerminals()
				return m, tea.Quit
			}
			cmd := m.doGlobal(action)
			m.applySize()
			if m.animating {
				return m, tea.Batch(cmd, tickCmd())
			}
			return m, cmd
		}
		if m.showHelp {
			m.showHelp = false
			return m, nil
		}
		if m.sbVisible && m.sbFocused {
			m.sidebarKey(ev)
			m.applySize()
			if m.animating {
				return m, tickCmd()
			}
			return m, nil
		}
		if m.e.Mode() == editor.InsertMode {
			m.km.HandleInsert(m.e, ev)
		} else {
			m.km.HandleView(m.e, ev)
		}
		m.applySize()
	}
	return m, nil
}

func (m *Model) doGlobal(action string) tea.Cmd {
	switch action {
	case keymap.Save:
		_ = m.e.Save()
	case keymap.Help:
		m.showHelp = !m.showHelp
	case keymap.Palette:
		m.openPalette()
	case keymap.Theme:
		m.showTheme = true
		m.themeSel = m.themeIdx
	case keymap.Sidebar:
		m.toggleSidebar()
	case keymap.SidebarSide:
		m.sbRight = !m.sbRight
		m.saveConfig()
	case keymap.TabNext:
		m.nextTab()
	case keymap.TabPrev:
		m.prevTab()
	case keymap.TabNew:
		m.newTab()
	case keymap.TabClose:
		m.requestClose()
	case keymap.CursorAddBelow:
		m.e.AddCaretBelow()
	case keymap.CursorAddAbove:
		m.e.AddCaretAbove()
	case keymap.Terminal:
		return m.toggleTerminal()
	}
	return nil
}

func (m Model) sidebarDisplayWidth() int {
	return int(math.Round(float64(m.sbWidth) * m.sbReveal))
}

func (m Model) editorOriginX() int {
	if m.sbRight {
		return 0
	}
	return m.sidebarDisplayWidth()
}

func (m *Model) applySize() {
	cols := m.width - m.sidebarDisplayWidth() - gutterWidth - 1
	rows := m.height - 2 - m.termHeight()
	if cols < 1 {
		cols = 1
	}
	if rows < 1 {
		rows = 1
	}
	m.e.SetViewport(rows, cols)
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	contentH := m.height - 2 - m.termHeight()
	if contentH < 1 {
		contentH = 1
	}
	sidebar := m.sidebarDisplayWidth()
	contentW := m.width - sidebar - gutterWidth - 1
	if contentW < 1 {
		contentW = 1
	}

	lexer := syntax.GetLexer(document.Extension(m.e.FileName()))
	offRow, offCol := m.e.Offset()
	cr, _ := m.e.Cursor()
	blank := lipgloss.NewStyle().Background(m.theme.Bg)

	lines := make([]string, 0, contentH)
	for r := 0; r < contentH; r++ {
		li := r + offRow
		if li >= m.e.LineCount() {
			lines = append(lines, m.st.gutter.Render(strings.Repeat(" ", gutterWidth))+blank.Width(contentW).Render(""))
			continue
		}
		cursorRow := li == cr && !m.sbFocused
		gutter := m.st.gutter
		if cursorRow {
			gutter = m.st.gutterHi
		}
		lines = append(lines, gutter.Render(fmt.Sprintf("%4d ", li+1))+m.renderLine(li, lexer, contentW, offCol, cursorRow))
	}
	editorPane := lipgloss.JoinHorizontal(lipgloss.Top, strings.Join(lines, "\n"), m.scrollbar(contentH, offRow, m.e.LineCount()))

	body := editorPane
	if sidebar > 0 {
		if m.sbRight {
			body = lipgloss.JoinHorizontal(lipgloss.Top, editorPane, m.sidebarView(contentH))
		} else {
			body = lipgloss.JoinHorizontal(lipgloss.Top, m.sidebarView(contentH), editorPane)
		}
	}
	parts := []string{m.titleBar(m.width), body}
	if m.termHeight() > 0 {
		parts = append(parts, m.terminalView(m.width, m.termHeight()))
	}
	parts = append(parts, m.statusView())
	screen := lipgloss.JoinVertical(lipgloss.Left, parts...)

	switch {
	case m.showPalette:
		return m.overlay(m.paletteView())
	case m.showTheme:
		return m.overlay(m.themePickerView())
	case m.showHelp:
		return m.overlay(m.helpView())
	}
	return screen
}

func (m Model) overlay(box string) string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box,
		lipgloss.WithWhitespaceChars(" "), lipgloss.WithWhitespaceForeground(m.theme.Bg))
}

type cellKey struct {
	fg     lipgloss.Color
	bg     lipgloss.Color
	cursor bool
}

type dcell struct {
	r   rune
	key cellKey
}

func (m Model) renderLine(row int, lexer chroma.Lexer, width, offCol int, cursorRow bool) string {
	line := m.e.Line(row)
	colors := lineColors(line, lexer, m.theme)
	caretCols := m.caretColsOnRow(row)
	selActive, sr, sc, er, ec := orderedSelection(m.e)

	lineBg := m.theme.Bg
	if cursorRow {
		lineBg = m.theme.CurrentLine
	}

	isCaret := func(col int) bool {
		for _, c := range caretCols {
			if c == col {
				return true
			}
		}
		return false
	}

	cells := make([]dcell, 0, len(line)+tabStop)
	screen := 0
	for i, r := range line {
		fg := colors[i]
		bg := lineBg
		if selActive && inSelection(row, i, sr, sc, er, ec) {
			bg = m.theme.Selection
		}
		cur := isCaret(i)
		if r == '\t' {
			n := tabStop - (screen % tabStop)
			for k := 0; k < n; k++ {
				cells = append(cells, dcell{' ', cellKey{fg, bg, cur && k == 0}})
				screen++
			}
			continue
		}
		cells = append(cells, dcell{r, cellKey{fg, bg, cur}})
		screen++
	}
	if isCaret(len(line)) {
		cells = append(cells, dcell{' ', cellKey{m.theme.Fg, lineBg, true}})
	}

	at := func(i int) dcell {
		if i >= 0 && i < len(cells) {
			return cells[i]
		}
		return dcell{' ', cellKey{m.theme.Fg, lineBg, false}}
	}

	var b strings.Builder
	p := 0
	for p < width {
		c := at(offCol + p)
		run := []rune{c.r}
		q := p + 1
		for q < width {
			c2 := at(offCol + q)
			if c2.key != c.key {
				break
			}
			run = append(run, c2.r)
			q++
		}
		b.WriteString(m.styleForKey(c.key).Render(string(run)))
		p = q
	}
	return b.String()
}

func (m Model) caretColsOnRow(row int) []int {
	if m.sbFocused || (m.termVisible && m.termFocused) {
		return nil
	}
	var cols []int
	for _, c := range m.e.Carets() {
		if c[0] == row {
			cols = append(cols, c[1])
		}
	}
	return cols
}

func (m Model) styleForKey(k cellKey) lipgloss.Style {
	if k.cursor {
		return m.st.cursor
	}
	return lipgloss.NewStyle().Foreground(k.fg).Background(k.bg)
}

func inSelection(row, col, sr, sc, er, ec int) bool {
	if row < sr || row > er {
		return false
	}
	if row == sr && col < sc {
		return false
	}
	if row == er && col >= ec {
		return false
	}
	return true
}

func orderedSelection(e *editor.Editor) (bool, int, int, int, int) {
	active, sr, sc, er, ec := e.Selection()
	if !active {
		return false, 0, 0, 0, 0
	}
	if sr > er || (sr == er && sc > ec) {
		sr, sc, er, ec = er, ec, sr, sc
	}
	return true, sr, sc, er, ec
}

func lineColors(line []rune, lexer chroma.Lexer, th theme.Theme) []lipgloss.Color {
	colors := make([]lipgloss.Color, len(line))
	for i := range colors {
		colors[i] = th.Fg
	}
	it, err := lexer.Tokenise(nil, string(line))
	if err != nil {
		return colors
	}
	i := 0
	for t := it(); t != chroma.EOF; t = it() {
		c := th.Token(t.Type)
		for range t.Value {
			if i < len(colors) {
				colors[i] = c
				i++
			}
		}
	}
	return colors
}

func (m Model) statusView() string {
	mode := " VIEW "
	if m.e.Mode() == editor.InsertMode {
		mode = " INSERT "
	}
	pill := m.st.modePill.Render(mode)
	arrow := m.st.plArrow.Render(icons.Powerline())

	name := baseName(m.e.FileName())
	label := m.e.FileName()
	if label == "" {
		label = "untitled"
	}
	if m.e.Modified() {
		label += " ●"
	}
	left := " " + icons.File(name) + " " + label + " "
	if msg, ok := m.e.ActiveMessage(); ok {
		left = " " + msg + " "
	}

	cr, cc := m.e.Cursor()
	info := fmt.Sprintf(" %s ", m.e.Language())
	if n := m.e.CaretCount(); n > 1 {
		info = fmt.Sprintf(" %d cursors · %s ", n, m.e.Language())
	}
	posArrow := lipgloss.NewStyle().Foreground(m.theme.Accent).Background(m.theme.StatusBg).Render(icons.PowerlineLeft())
	posPill := lipgloss.NewStyle().Foreground(m.theme.CursorFg).Background(m.theme.Accent).Bold(true).Render(fmt.Sprintf(" %d:%d ", cr+1, cc+1))

	head := lipgloss.Width(pill) + lipgloss.Width(arrow)
	tail := lipgloss.Width(info) + lipgloss.Width(posArrow) + lipgloss.Width(posPill)
	fill := m.width - head - lipgloss.Width(left) - tail
	if fill < 0 {
		info = ""
		tail = lipgloss.Width(posArrow) + lipgloss.Width(posPill)
		fill = m.width - head - lipgloss.Width(left) - tail
	}
	if fill < 0 {
		mid := m.st.statusBar.Render(fitLeft(left, m.width-head))
		return pill + arrow + mid
	}
	mid := m.st.statusBar.Render(left + strings.Repeat(" ", fill) + info)
	return pill + arrow + mid + posArrow + posPill
}

func fitLeft(s string, width int) string {
	if width < 0 {
		width = 0
	}
	if w := lipgloss.Width(s); w < width {
		return s + strings.Repeat(" ", width-w)
	}
	return ansi.Truncate(s, width, "")
}

func baseName(path string) string {
	if path == "" {
		return "untitled"
	}
	if i := strings.LastIndexByte(path, '/'); i >= 0 {
		return path[i+1:]
	}
	return path
}
