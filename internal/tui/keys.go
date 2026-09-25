package tui

import (
	"github.com/arthurlch/cub/internal/input"
	tea "github.com/charmbracelet/bubbletea"
)

func fromTea(msg tea.KeyMsg) input.Event {
	alt := msg.Alt

	switch msg.Type {
	case tea.KeyRunes:
		if len(msg.Runes) > 0 {
			return input.Event{Key: input.KeyRune, Rune: msg.Runes[0], Alt: alt}
		}
		return input.Event{}
	case tea.KeySpace:
		return input.Event{Key: input.KeyRune, Rune: ' ', Alt: alt}
	case tea.KeyEnter:
		return input.Event{Key: input.KeyEnter, Alt: alt}
	case tea.KeyTab:
		return input.Event{Key: input.KeyTab, Alt: alt}
	case tea.KeyEsc:
		return input.Event{Key: input.KeyEsc}
	case tea.KeyBackspace:
		return input.Event{Key: input.KeyBackspace}
	case tea.KeyDelete:
		return input.Event{Key: input.KeyDelete}
	case tea.KeyUp:
		return input.Event{Key: input.KeyUp, Alt: alt}
	case tea.KeyDown:
		return input.Event{Key: input.KeyDown, Alt: alt}
	case tea.KeyLeft:
		return input.Event{Key: input.KeyLeft, Alt: alt}
	case tea.KeyRight:
		return input.Event{Key: input.KeyRight, Alt: alt}
	case tea.KeyHome:
		return input.Event{Key: input.KeyHome, Alt: alt}
	case tea.KeyEnd:
		return input.Event{Key: input.KeyEnd, Alt: alt}
	case tea.KeyPgUp:
		return input.Event{Key: input.KeyPgUp, Alt: alt}
	case tea.KeyPgDown:
		return input.Event{Key: input.KeyPgDn, Alt: alt}
	case tea.KeyCtrlUp:
		return input.Event{Key: input.KeyUp, Ctrl: true}
	case tea.KeyCtrlDown:
		return input.Event{Key: input.KeyDown, Ctrl: true}
	case tea.KeyCtrlLeft:
		return input.Event{Key: input.KeyLeft, Ctrl: true}
	case tea.KeyCtrlRight:
		return input.Event{Key: input.KeyRight, Ctrl: true}
	}

	if msg.Type >= tea.KeyCtrlA && msg.Type <= tea.KeyCtrlZ {
		return input.Event{Key: input.KeyRune, Rune: rune('a' + int(msg.Type-tea.KeyCtrlA)), Ctrl: true}
	}
	return input.Event{}
}
