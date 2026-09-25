package keymap

import (
	"strconv"
	"strings"

	"github.com/arthurlch/cub/internal/editor"
	"github.com/arthurlch/cub/internal/input"
)

// NOTE: I have an emotionnal meltdown about this file.
// I don't like it. It's a mess.
// But it works. And I don't have the energy to refactor it right now.
// I also have another emotional meltdown about the fact that the UX atm isnt great
// but it works for me !
const (
	Quit        = "quit"
	Save        = "save"
	Help        = "help"
	Palette     = "palette"
	Theme       = "theme"
	Sidebar     = "sidebar"
	SidebarSide = "sidebar_side"
	TabNext     = "tab_next"
	TabPrev     = "tab_prev"
	TabNew      = "tab_new"
	TabClose    = "tab_close"

	CursorAddBelow = "cursor_add_below"
	CursorAddAbove = "cursor_add_above"

	Terminal = "terminal"

	Insert        = "insert"
	Append        = "append"
	OpenBelow     = "open_below"
	OpenAbove     = "open_above"
	MoveLeft      = "move_left"
	MoveDown      = "move_down"
	MoveUp        = "move_up"
	MoveRight     = "move_right"
	LineStart     = "line_start"
	FirstNonBlank = "first_nonblank"
	LineEnd       = "line_end"
	WordNext      = "word_next"
	WordPrev      = "word_prev"
	GotoTop       = "goto_top"
	GotoBottom    = "goto_bottom"
	DeleteChar    = "delete_char"
	DeleteLine    = "delete_line"
	YankLine      = "yank_line"
	Paste         = "paste"
	PasteBefore   = "paste_before"
	Visual        = "visual"
	Undo          = "undo"
	Redo          = "redo"
)

var globalActions = map[string]bool{
	Quit: true, Save: true, Help: true, Palette: true,
	Theme: true, Sidebar: true, SidebarSide: true,
	TabNext: true, TabPrev: true, TabNew: true, TabClose: true,
	CursorAddBelow: true, CursorAddAbove: true,
	Terminal: true,
}

var multiCursorSafe = map[string]bool{
	MoveLeft: true, MoveRight: true, MoveUp: true, MoveDown: true,
	LineStart: true, FirstNonBlank: true, LineEnd: true,
	Insert: true, Append: true,
}

var visualOperator = map[rune]bool{'y': true, 'd': true, 'x': true, 'p': true, 'v': true}

func DefaultChords() map[string]string {
	return map[string]string{
		Quit:        "ctrl+q",
		Save:        "ctrl+s",
		Help:        "ctrl+h",
		Palette:     "ctrl+p",
		Theme:       "ctrl+t",
		Sidebar:     "ctrl+b",
		SidebarSide: "ctrl+e",
		TabNext:     "ctrl+right",
		TabPrev:     "ctrl+left",
		TabNew:      "ctrl+n",
		TabClose:    "ctrl+w",

		CursorAddBelow: "ctrl+down",
		CursorAddAbove: "ctrl+up",

		Terminal: "ctrl+g",

		Insert:        "i",
		Append:        "a",
		OpenBelow:     "o",
		OpenAbove:     "O",
		MoveLeft:      "h",
		MoveDown:      "j",
		MoveUp:        "k",
		MoveRight:     "l",
		LineStart:     "0",
		FirstNonBlank: "^",
		LineEnd:       "$",
		WordNext:      "w",
		WordPrev:      "b",
		GotoTop:       "gg",
		GotoBottom:    "G",
		DeleteChar:    "x",
		DeleteLine:    "dd",
		YankLine:      "yy",
		Paste:         "p",
		PasteBefore:   "P",
		Visual:        "v",
		Undo:          "u",
		Redo:          "ctrl+r",
	}
}

func Parse(chord string) (input.Event, bool) {
	parts := strings.Split(strings.ToLower(strings.TrimSpace(chord)), "+")
	ev := input.Event{}
	key := parts[len(parts)-1]
	for _, p := range parts[:len(parts)-1] {
		switch p {
		case "ctrl":
			ev.Ctrl = true
		case "alt":
			ev.Alt = true
		default:
			return input.Event{}, false
		}
	}

	switch key {
	case "left":
		ev.Key = input.KeyLeft
	case "right":
		ev.Key = input.KeyRight
	case "up":
		ev.Key = input.KeyUp
	case "down":
		ev.Key = input.KeyDown
	case "home":
		ev.Key = input.KeyHome
	case "end":
		ev.Key = input.KeyEnd
	case "pgup", "pageup":
		ev.Key = input.KeyPgUp
	case "pgdn", "pagedown":
		ev.Key = input.KeyPgDn
	case "enter", "return":
		ev.Key = input.KeyEnter
	case "tab":
		ev.Key = input.KeyTab
	case "esc", "escape":
		ev.Key = input.KeyEsc
	case "space":
		ev.Key = input.KeyRune
		ev.Rune = ' '
	case "backspace":
		ev.Key = input.KeyBackspace
	case "delete", "del":
		ev.Key = input.KeyDelete
	default:
		r := []rune(chord)
		if len(parts) == 1 && len([]rune(parts[0])) == 1 {
			ev.Key = input.KeyRune
			ev.Rune = r[len(r)-1]
		} else if len([]rune(key)) == 1 {
			ev.Key = input.KeyRune
			ev.Rune = []rune(key)[0]
		} else {
			return input.Event{}, false
		}
	}
	return ev, true
}

type Bindings struct {
	global map[input.Event]string
	view   map[input.Event]string
	seq    map[input.Event]map[input.Event]string
}

func Build(chords map[string]string) Bindings {
	b := Bindings{
		global: map[input.Event]string{},
		view:   map[input.Event]string{},
		seq:    map[input.Event]map[input.Event]string{},
	}
	for action, chord := range chords {
		if seq, ok := asSequence(chord); ok {
			first := input.Event{Key: input.KeyRune, Rune: seq[0]}
			second := input.Event{Key: input.KeyRune, Rune: seq[1]}
			if b.seq[first] == nil {
				b.seq[first] = map[input.Event]string{}
			}
			b.seq[first][second] = action
			continue
		}
		ev, ok := Parse(chord)
		if !ok {
			continue
		}
		if globalActions[action] {
			b.global[ev] = action
		} else {
			b.view[ev] = action
		}
	}
	return b
}

func asSequence(chord string) ([]rune, bool) {
	if strings.Contains(chord, "+") {
		return nil, false
	}
	r := []rune(chord)
	if len(r) == 2 {
		return r, true
	}
	return nil, false
}

type Keymap struct {
	binds      Bindings
	pending    input.Event
	hasPending bool
	digits     string
}

func Default() *Keymap { return New(DefaultChords()) }

func New(chords map[string]string) *Keymap {
	return &Keymap{binds: Build(chords)}
}

func (k *Keymap) Global(ev input.Event) string { return k.binds.global[ev] }

func (k *Keymap) HandleInsert(e *editor.Editor, ev input.Event) {
	switch ev.Key {
	case input.KeyEnter:
		e.InsertNewline()
	case input.KeyTab:
		e.InsertTab()
	case input.KeyBackspace:
		e.Backspace()
	case input.KeyDelete:
		e.DeleteForward()
	case input.KeyEsc:
		e.ClearCarets()
		e.EnterView()
	case input.KeyUp:
		e.MoveUp()
	case input.KeyDown:
		e.MoveDown()
	case input.KeyLeft:
		e.MoveLeft()
	case input.KeyRight:
		e.MoveRight()
	case input.KeyHome:
		e.Home()
	case input.KeyEnd:
		e.LineEnd()
	case input.KeyPgUp:
		e.PageUp()
	case input.KeyPgDn:
		e.PageDown()
	case input.KeyRune:
		if ev.Rune != 0 && !ev.Ctrl && !ev.Alt {
			e.InsertRune(ev.Rune)
		}
	}
}

func (k *Keymap) HandleView(e *editor.Editor, ev input.Event) {
	switch ev.Key {
	case input.KeyUp:
		e.MoveUp()
		k.reset()
		return
	case input.KeyDown:
		e.MoveDown()
		k.reset()
		return
	case input.KeyLeft:
		e.MoveLeft()
		k.reset()
		return
	case input.KeyRight:
		e.MoveRight()
		k.reset()
		return
	case input.KeyHome:
		e.Home()
		k.reset()
		return
	case input.KeyEnd:
		e.LineEnd()
		k.reset()
		return
	case input.KeyPgUp:
		e.PageUp()
		k.reset()
		return
	case input.KeyPgDn:
		e.PageDown()
		k.reset()
		return
	case input.KeyDelete, input.KeyBackspace:
		e.DeleteSelection()
		k.reset()
		return
	case input.KeyEsc:
		e.EndSelection()
		e.ClearCarets()
		k.reset()
		return
	}

	if e.SelectionActive() && ev.Key == input.KeyRune && !ev.Ctrl && !ev.Alt && visualOperator[ev.Rune] {
		k.reset()
		switch ev.Rune {
		case 'y':
			e.Copy()
			e.EndSelection()
		case 'd', 'x':
			e.Cut()
		case 'p':
			e.Paste()
		case 'v':
			e.EndSelection()
		}
		return
	}

	if k.hasPending {
		if action, ok := k.binds.seq[k.pending][ev]; ok {
			k.reset()
			k.doView(e, action)
			return
		}
		k.hasPending = false
	}
	if _, ok := k.binds.seq[ev]; ok {
		k.pending = ev
		k.hasPending = true
		return
	}

	if ev.Key == input.KeyRune && ev.Rune >= '0' && ev.Rune <= '9' && !ev.Ctrl && !ev.Alt {
		k.digits += string(ev.Rune)
		return
	}

	if action, ok := k.binds.view[ev]; ok {
		k.doView(e, action)
		return
	}
	k.digits = ""
}

func (k *Keymap) doView(e *editor.Editor, action string) {
	if !multiCursorSafe[action] {
		e.ClearCarets()
	}
	switch action {
	case Insert:
		e.EnterInsert()
	case Append:
		e.Append()
	case OpenBelow:
		e.OpenBelow()
	case OpenAbove:
		e.OpenAbove()
	case MoveLeft:
		e.MoveLeft()
	case MoveDown:
		e.MoveDown()
	case MoveUp:
		e.MoveUp()
	case MoveRight:
		e.MoveRight()
	case LineStart:
		e.Home()
	case FirstNonBlank:
		e.LineStart()
	case LineEnd:
		e.LineEnd()
	case WordNext:
		e.WordNext()
	case WordPrev:
		e.WordPrev()
	case GotoTop:
		e.GotoTop()
	case GotoBottom:
		if k.digits != "" {
			n, _ := strconv.Atoi(k.digits)
			e.GotoLine(n)
		} else {
			e.GotoBottom()
		}
	case DeleteChar:
		e.DeleteChar()
	case DeleteLine:
		e.DeleteLine()
	case YankLine:
		e.YankLine()
	case Paste:
		e.Paste()
	case PasteBefore:
		e.PasteBefore()
	case Visual:
		if e.SelectionActive() {
			e.EndSelection()
		} else {
			e.StartSelection()
		}
	case Undo:
		e.Undo()
	case Redo:
		e.Redo()
	}
	k.digits = ""
}

func (k *Keymap) reset() {
	k.hasPending = false
	k.digits = ""
}
