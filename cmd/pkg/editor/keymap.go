package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/nsf/termbox-go"
)

type KeyAction func(st *state.State, keyEvent termbox.Event)

var viewKeyBindings = map[termbox.Key]KeyAction{
	termbox.KeyArrowUp:    handleNavigation,
	termbox.KeyArrowDown:  handleNavigation,
	termbox.KeyArrowLeft:  handleNavigation,
	termbox.KeyArrowRight: handleNavigation,
	termbox.KeyHome:       handleNavigation,
	termbox.KeyEnd:        handleNavigation,
	termbox.KeyPgup:       handleNavigation,
	termbox.KeyPgdn:       handleNavigation,
}

var viewRuneBindings = map[rune]KeyAction{
	's': func(st *state.State, _ termbox.Event) { StartSelection(st) },
	'z': func(st *state.State, _ termbox.Event) { EndSelection(st) },
	'c': func(st *state.State, _ termbox.Event) { CopySelection(st); EndSelection(st) },
	'x': func(st *state.State, _ termbox.Event) { CutSelection(st); EndSelection(st) },
	'v': func(st *state.State, _ termbox.Event) { PasteSelection(st) },

	'k': func(st *state.State, _ termbox.Event) { moveUp(st) },
	'j': func(st *state.State, _ termbox.Event) { moveRight(st) },
	'i': func(st *state.State, _ termbox.Event) { moveLeft(st) },
	'm': func(st *state.State, _ termbox.Event) { moveDown(st) },
	'w': func(st *state.State, _ termbox.Event) { moveToNextWord(st) },
	'b': func(st *state.State, _ termbox.Event) { moveToPreviousWord(st) },
	'(': func(st *state.State, _ termbox.Event) { moveToMatchingBracket(st, '(') },
	')': func(st *state.State, _ termbox.Event) { moveToMatchingBracket(st, ')') },
	'e': func(st *state.State, _ termbox.Event) { moveToNextEmptyLine(st) },
	'E': func(st *state.State, _ termbox.Event) { moveToPreviousEmptyLine(st) },
	'^': func(st *state.State, _ termbox.Event) { moveToLineStart(st) },
	'$': func(st *state.State, _ termbox.Event) { moveToLineEnd(st) },
	'g': func(st *state.State, _ termbox.Event) { st.MoveCursor(0, 0) },
	'G': func(st *state.State, _ termbox.Event) { jumpToLine(st); st.LineNumberBuffer = "" },
}

func init() {
	for _, d := range "0123456789" {
		digit := d
		viewRuneBindings[digit] = func(st *state.State, _ termbox.Event) {
			st.LineNumberBuffer += string(digit)
		}
	}
}

func isDeleteKey(k termbox.Key) bool {
	return k == termbox.KeyDelete || k == termbox.KeyBackspace || k == termbox.KeyBackspace2
}

func tryViewKeyBinding(st *state.State, keyEvent termbox.Event) bool {
	if action, ok := viewKeyBindings[keyEvent.Key]; ok {
		action(st, keyEvent)
		return true
	}
	return false
}

func tryPendingSequence(st *state.State, keyEvent termbox.Event) bool {
	switch keyEvent.Ch {
	case 'd':
		if st.LastKey == 'd' {
			DeleteCurrentLine(st)
			st.LastKey = 0
		} else {
			st.LastKey = 'd'
		}
		return true
	case 'a':
		if st.LastKey == 'a' {
			SelectAll(st)
			st.LastKey = 0
		} else {
			st.LastKey = 'a'
		}
		return true
	}
	return false
}

func tryViewRuneBinding(st *state.State, keyEvent termbox.Event) bool {
	if action, ok := viewRuneBindings[keyEvent.Ch]; ok {
		st.LastKey = 0
		action(st, keyEvent)
		return true
	}
	return false
}
