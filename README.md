# -- Cub Text Editor --

<p align="center">
  <img src="https://i.ibb.co/88MWThZ/d9f9d3aa-2ae6-47d9-96bc-2735eda584f9.webp" alt="Cub Logo" width="300" height="300">
</p>

![GitHub release (latest by date)](https://img.shields.io/github/v/release/arthurlch/cub)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/arthurlch/cub/build_and_test.yml)
![GitHub issues](https://img.shields.io/github/issues/arthurlch/cub)

## Overview

**Cub** is a fast, modern, terminal-based text editor. It pairs a modal, Vim/Kakoune-style editing model with a polished TUI — a file sidebar, tabs, a fuzzy file finder, mouse support, and truecolor themes — built on the [Charm](https://charm.sh) stack (Bubble Tea + Lip Gloss).

Cub has two modes:

- **View mode** — navigate, select, and manipulate text.
- **Insert mode** — type and edit text (`i` to enter, `Esc` to leave).

### Features

- **Modal editing** — vim-style `hjkl` motions, `w`/`b` words, `gg`/`G`, visual mode (`v`), and line registers (`yy`/`dd`/`p`).
- **Multiple cursors** — stack cursors across lines and edit them all at once, with grouped undo.
- **Integrated terminal** — one or more real shells in a bottom panel, powered by a built-in VT emulator.
- **250+ languages** and **24 truecolor colorschemes** (dark & light) via [chroma](https://github.com/alecthomas/chroma), with real per-token syntax highlighting.
- **Tabs** — open many files at once, each with its own cursor, selection, and undo history.
- **File sidebar** — a toggleable, resizable file explorer.
- **Command palette** — fuzzy file finder (`Ctrl+P`).
- **Configurable** — rebind every global and view-mode key, and your theme and sidebar side are remembered, all in one JSON config.
- **Full mouse support** — click to place the cursor, drag to select, drag the border to resize, wheel to scroll, click tabs.
- **Fast core** — edits and undo are O(change), not O(file), so it stays instant on large files.

---

## Keybindings

Every key below is a default — global and view-mode bindings are fully rebindable from the config file. See [Configuration](#configuration).

### Global (any mode)

| Key | Action |
| --- | --- |
| `Ctrl+S` | Save |
| `Ctrl+Q` | Quit |
| `Ctrl+H` | Toggle help overlay |
| `Ctrl+P` | Command palette (fuzzy file finder) |
| `Ctrl+T` | Theme switcher |
| `Ctrl+B` | File sidebar (show / focus / return focus) |
| `Ctrl+E` | Move the sidebar to the other side (left / right) |
| `Ctrl+←` / `Ctrl+→` | Previous / next tab |
| `Ctrl+N` | New tab |
| `Ctrl+W` | Close tab |
| `Ctrl+↑` / `Ctrl+↓` | Add a cursor above / below |
| `Ctrl+G` | Toggle the integrated terminal |

### View mode (vim-style)

| Key | Action |
| --- | --- |
| `i` / `a` | Insert before / append after the cursor |
| `o` / `O` | Open a new line below / above and insert |
| `h` `j` `k` `l` or arrows | Move left / down / up / right |
| `0` / `^` / `$` | Line start / first non-blank / line end |
| `w` / `b` | Next / previous word |
| `gg` / `G` | Top / bottom of file (`<number>G` jumps to a line, e.g. `42G`) |
| `PgUp` / `PgDn` | Scroll half a screen |
| `x` | Delete the character under the cursor |
| `dd` / `yy` | Delete / yank the current line |
| `p` / `P` | Paste after / before |
| `u` / `Ctrl+R` | Undo / redo |

### Visual (selection)

Press `v` to start a selection, then move (`hjkl`, word, line motions) to extend it.

| Key | Action |
| --- | --- |
| `v` | Start / cancel visual selection |
| `y` | Yank (copy) the selection |
| `d` / `x` | Delete (cut) the selection |
| `p` | Paste over the selection |
| `Esc` | Cancel the selection |

### Insert mode

| Key | Action |
| --- | --- |
| `Esc` | Return to view mode |
| any character | Insert it |
| `Enter` | New line |
| `Tab` | Insert 4 spaces |
| `Backspace` / `Delete` | Delete before / at the cursor |
| arrows, `Home`/`End`, `PgUp`/`PgDn` | Move the cursor |

---

## Tabs

Open multiple files at once, shown as a clickable tab strip along the top. Each tab is fully independent — its own cursor, selection, viewport, and undo history.

- Opening a file (sidebar or command palette) opens it in a tab, or focuses its existing tab if it's already open. An empty, unsaved scratch buffer is reused instead of adding a tab.
- **`Ctrl+N`** — open a new empty tab.
- **`Ctrl+→` / `Ctrl+←`** — next / previous tab (wraps around).
- **`Ctrl+W`** — close the current tab. If it has unsaved changes, press again to discard. Closing the last tab leaves a fresh empty buffer.
- **Click** a tab to switch, or click its **`✕`** to close it (a dirty tab warns first, like `Ctrl+W`). The strip scrolls to keep the active tab visible when there are more tabs than fit. `●` marks unsaved changes.

## Multiple cursors

- **`Ctrl+↓`** — add a cursor on the line below; **`Ctrl+↑`** — add one on the line above. Repeat to stack more.
- Enter insert mode (`i`) and type: every cursor edits at once. Motions (`hjkl`, arrows, `Home`/`End`) move all cursors together.
- Backspace, delete, and newline apply at every cursor, and a single **`u`** undoes the whole multi-cursor edit as one step.
- **`Esc`**, a click, or any jump (word/goto/page) collapses back to a single cursor. The status bar shows the active cursor count.

## Integrated terminal

- **`Ctrl+G`** — toggle a real shell (`$SHELL`) in a panel across the bottom of the editor. It runs through a built-in VT emulator, so colors, prompts, and full-screen programs work.
- **Focus** — `Ctrl+G` is smart: it shows + focuses the panel, focuses it again if you'd clicked away, or hides it when it's already focused. While the terminal is focused, keystrokes go straight to the shell; **`Esc`** hands focus back to the editor without closing the panel. Click either pane to focus it. `Ctrl+Q` always quits cub.
- **Multiple terminals** — while focused, **`Ctrl+N`** opens another shell, **`Ctrl+W`** closes the current one, and **`Ctrl+←` / `Ctrl+→`** switch between them. The header shows numbered tabs (`1 ✕  2 ✕  +`); click a number to switch, its **`✕`** to close it, or **`+`** to add one.
- **Resize it** — while focused, **`Ctrl+↑` / `Ctrl+↓`** grow/shrink the panel, or **drag its header bar** with the mouse.
- Each panel remembers its size, resizes with the window, and every shell is cleaned up on exit.

## File sidebar

- **`Ctrl+B`** — show and focus the sidebar. If it's already open, `Ctrl+B` flips focus between the sidebar and the editor. The focused pane is highlighted (accent border/selection); the other is dimmed.
- **`Ctrl+E`** — dock the sidebar on the other side (left ↔ right).
- **`Esc`** — close the sidebar (returns focus to the editor).
- **`j` / `k`** or **↑ / ↓** — move the selection.
- **`Enter`** or **`l`** — open the selected file, or expand/collapse the selected folder.
- **`h`** — collapse the current folder.
- **`<` / `>`** (also **`-` / `+`**) — resize the sidebar; text reflows live.

## Command palette

- **`Ctrl+P`** — fuzzy file finder. Type to filter files in the project, **↑ / ↓** to select, **`Enter`** to open, **`Esc`** to cancel.

## Themes & languages

Cub highlights **250+ languages** and ships **24 truecolor colorschemes** (dark and light), all powered by [chroma](https://github.com/alecthomas/chroma). A colorscheme drives per-token syntax colors *and* the derived UI colors (cursor line, status bar, sidebar, borders). Schemes include catppuccin (mocha/macchiato/latte), dracula, tokyonight (night/storm/day), nord, gruvbox (dark/light), monokai, onedark, github-dark, rose-pine (main/moon/dawn), solarized (dark/light), doom-one, xcode, and more.

- **`Ctrl+T`** — open the theme switcher. **↑ / ↓** previews live, **`Enter`** applies, **`Esc`** cancels and reverts.

## Configuration

Cub reads a JSON config from `~/.config/cub/config.json` (or `$XDG_CONFIG_HOME/cub/config.json`). It is written with the defaults on first run, so you can open it and edit it right away.

```json
{
  "theme": "catppuccin-mocha",
  "sidebar_right": false,
  "keys": {
    "quit": "ctrl+q",
    "save": "ctrl+s",
    "insert": "i",
    "move_left": "h",
    "delete_line": "dd"
  }
}
```

- **`theme`** and **`sidebar_right`** are saved automatically whenever you switch theme (`Ctrl+T`) or flip the sidebar side (`Ctrl+E`).
- **`keys`** maps an action to a chord. Every global and view-mode action is rebindable; insert-mode typing is fixed. Any action you leave out keeps its default.
- cub also writes a `version` field it manages for config migrations — leave it alone; when defaults change across versions it refreshes your keymap while keeping your theme and sidebar side.

Chord syntax: a single character (`i`, `$`), a named key (`enter`, `esc`, `tab`, `space`, `backspace`, `delete`, `home`, `end`, `pgup`, `pgdn`, arrow names `left`/`right`/`up`/`down`), with optional `ctrl+` / `alt+` prefixes (e.g. `ctrl+s`, `ctrl+right`). A two-character value is a sequence, like the default `dd` or `gg`.

| Action | Default | Action | Default |
| --- | --- | --- | --- |
| `quit` | `ctrl+q` | `insert` / `append` | `i` / `a` |
| `save` | `ctrl+s` | `open_below` / `open_above` | `o` / `O` |
| `help` | `ctrl+h` | `move_left` / `move_down` | `h` / `j` |
| `palette` | `ctrl+p` | `move_up` / `move_right` | `k` / `l` |
| `theme` | `ctrl+t` | `line_start` / `first_nonblank` / `line_end` | `0` / `^` / `$` |
| `sidebar` | `ctrl+b` | `word_next` / `word_prev` | `w` / `b` |
| `sidebar_side` | `ctrl+e` | `goto_top` / `goto_bottom` | `gg` / `G` |
| `tab_next` / `tab_prev` | `ctrl+right` / `ctrl+left` | `delete_char` / `delete_line` | `x` / `dd` |
| `tab_new` / `tab_close` | `ctrl+n` / `ctrl+w` | `yank_line` / `paste` / `paste_before` | `yy` / `p` / `P` |
| `cursor_add_below` / `cursor_add_above` | `ctrl+down` / `ctrl+up` | `visual` | `v` |
| `terminal` | `ctrl+g` | `undo` / `redo` | `u` / `ctrl+r` |

In visual mode (`v`), `y` yanks, `d`/`x` delete, and `p` pastes over the selection — these are fixed, not rebindable.

## Mouse

- **Click** in the editor to place the cursor; **drag** to select text.
- **Click** a tab to switch, or its **`✕`** to close it; **click** a sidebar entry to open a file or toggle a folder.
- In the terminal header, **click** a numbered tab to switch, its **`✕`** to close it, or **`+`** to open one; **click** the panel body to focus it.
- **Drag the sidebar border** (`│`) to resize it; **drag the terminal header** to resize the terminal.
- **Scroll wheel** to move through the buffer, or the file tree when hovering the sidebar.

---

## Installation

### Homebrew

```bash
brew tap arthurlch/cub
brew install cub
```

### From source

```bash
git clone https://github.com/arthurlch/cub.git
cd cub
make build      # builds ./cub
```

Requires Go 1.24+.

---

## Usage

```bash
cub [file]      # open a file (created on save if it doesn't exist)
cub             # start with an empty buffer
cub --help      # show keybindings and exit
cub --version   # print version and exit
```

The UI uses Nerd Font glyphs (file-type icons and powerline separators). If your terminal isn't using a Nerd Font, run with `CUB_ASCII=1` to fall back to plain characters:

```bash
CUB_ASCII=1 cub main.go
```

Truecolor themes look best in a truecolor-capable terminal (kitty, WezTerm, Alacritty, Ghostty, iTerm2, etc.).

---

## Tips

- Enter insert mode with `i` (or `a` to append, `o`/`O` to open a line), get back to view mode with `Esc`.
- Select with `v` (visual), move to extend, then `y` (yank) / `d` (delete), and `p` to paste.
- Duplicate a line with `yy` then `p`; delete one with `dd`.
- Jump to a line with `<number>G` (e.g. `120G`); jump to the top with `gg`, the bottom with `G`.
- Undo/redo with `u` / `Ctrl+R`.
- Drop extra cursors with `Ctrl+↑` / `Ctrl+↓`, then type once to edit every line.
- Pop a shell with `Ctrl+G`; `Esc` jumps back to the editor.
- Open files fast with the fuzzy finder (`Ctrl+P`) or the sidebar (`Ctrl+B`).
- Press `Ctrl+H` any time for the in-app keybinding overlay.

---
