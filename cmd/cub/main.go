package main

import (
	"fmt"
	"os"

	"github.com/arthurlch/cub/internal/config"
	"github.com/arthurlch/cub/internal/editor"
	"github.com/arthurlch/cub/internal/filetree"
	"github.com/arthurlch/cub/internal/keymap"
	"github.com/arthurlch/cub/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

const sidebarWidth = 26

var version = "dev"

const usage = `cub - a fast, pretty terminal text editor

usage:
  cub [file]

flags:
  -h, --help      show this help and exit
  -v, --version   show version and exit

keys (view mode):
  i a            insert / append          o / O          open line below / above
  Esc            back to view mode        u / Ctrl+R     undo / redo
  h j k l        move                     w / b          word forward / back
  0 ^ $          start / first / end      gg / G         top / bottom (<n>G line)
  v              visual select            y d x          yank / delete / cut
  dd yy          delete / yank line       p / P          paste after / before

keys (global):
  Ctrl+S         save                     Ctrl+Q         quit
  Ctrl+B         file sidebar             Ctrl+P         fuzzy file finder
  Ctrl+T         switch theme             Ctrl+G         toggle terminal
  Ctrl+←/→       prev / next tab          Ctrl+W         close tab
  Ctrl+N         new tab                  Ctrl+E         sidebar left / right
  Ctrl+↑/↓       add cursor above/below   Ctrl+H         help

config:
  ~/.config/cub/config.json   theme, sidebar side, and all keybindings
                              (respects XDG_CONFIG_HOME; written on first run)

env:
  CUB_ASCII=1    use plain glyphs when the terminal has no Nerd Font
`

func run() {
	e := editor.New()
	if len(os.Args) > 1 {
		if err := e.Open(os.Args[1]); err != nil {
			e.ShowMessage(fmt.Sprintf("open failed: %v", err))
		}
	}

	tree, _ := filetree.New(".")

	cfg := config.Load()
	km := keymap.New(cfg.Chords())
	model := tui.New(e, tree, sidebarWidth).WithConfig(cfg, km)

	program := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := program.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version":
			fmt.Println("cub", version)
			return
		case "-h", "--help":
			fmt.Print(usage)
			return
		}
	}
	run()
}
