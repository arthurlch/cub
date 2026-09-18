package input

import "github.com/nsf/termbox-go"

func PollEvent() termbox.Event {
	return termbox.PollEvent()
}
