package theme

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Name  string
	style *chroma.Style

	Bg          lipgloss.Color
	Fg          lipgloss.Color
	CurrentLine lipgloss.Color
	Gutter      lipgloss.Color
	GutterHi    lipgloss.Color
	Selection   lipgloss.Color
	CursorBg    lipgloss.Color
	CursorFg    lipgloss.Color
	Accent      lipgloss.Color
	StatusBg    lipgloss.Color
	StatusFg    lipgloss.Color
	SidebarFg   lipgloss.Color
	SidebarDir  lipgloss.Color
	Border      lipgloss.Color
}

var schemeNames = []string{
	"catppuccin-mocha",
	"catppuccin-macchiato",
	"dracula",
	"tokyonight-night",
	"tokyonight-storm",
	"nord",
	"gruvbox",
	"monokai",
	"onedark",
	"github-dark",
	"rose-pine",
	"rose-pine-moon",
	"solarized-dark",
	"doom-one",
	"native",
	"paraiso-dark",
	"xcode-dark",
	"catppuccin-latte",
	"rose-pine-dawn",
	"tokyonight-day",
	"gruvbox-light",
	"solarized-light",
	"xcode",
	"modus-operandi",
}

var Themes []Theme

func init() {
	for _, n := range schemeNames {
		Themes = append(Themes, fromStyle(n))
	}
}

func Default() Theme { return Themes[0] }

func (t Theme) Token(tt chroma.TokenType) lipgloss.Color {
	e := t.style.Get(tt)
	if e.Colour.IsSet() {
		return lipgloss.Color(e.Colour.String())
	}
	return t.Fg
}

func fromStyle(name string) Theme {
	s := styles.Get(name)
	base := s.Get(chroma.Background)

	bg := col(base.Background, lipgloss.Color("#1e1e2e"))
	fg := col(base.Colour, lipgloss.Color("#e6e6e6"))
	kw := col(s.Get(chroma.Keyword).Colour, fg)
	fn := col(s.Get(chroma.NameFunction).Colour, kw)
	comment := col(s.Get(chroma.Comment).Colour, lerp(fg, bg, 0.5))
	gutter := col(s.Get(chroma.LineNumbers).Colour, comment)
	curLine := col(s.Get(chroma.LineHighlight).Background, lerp(bg, fg, 0.08))

	return Theme{
		Name: name, style: s,
		Bg: bg, Fg: fg, CurrentLine: curLine,
		Gutter: gutter, GutterHi: kw,
		Selection: lerp(bg, kw, 0.3),
		CursorBg:  kw, CursorFg: bg,
		Accent:   kw,
		StatusBg: lerp(bg, fg, 0.14), StatusFg: fg,
		SidebarFg: lerp(fg, bg, 0.15), SidebarDir: fn,
		Border: kw,
	}
}

func col(c chroma.Colour, fallback lipgloss.Color) lipgloss.Color {
	if c.IsSet() {
		return lipgloss.Color(c.String())
	}
	return fallback
}

func lerp(a, b lipgloss.Color, t float64) lipgloss.Color {
	ar, ag, ab := hex(a)
	br, bg, bb := hex(b)
	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x",
		int(float64(ar)+(float64(br)-float64(ar))*t),
		int(float64(ag)+(float64(bg)-float64(ag))*t),
		int(float64(ab)+(float64(bb)-float64(ab))*t),
	))
}

func hex(c lipgloss.Color) (r, g, b int) {
	s := strings.TrimPrefix(string(c), "#")
	if len(s) != 6 {
		return 128, 128, 128
	}
	ri, _ := strconv.ParseInt(s[0:2], 16, 0)
	gi, _ := strconv.ParseInt(s[2:4], 16, 0)
	bi, _ := strconv.ParseInt(s[4:6], 16, 0)
	return int(ri), int(gi), int(bi)
}
