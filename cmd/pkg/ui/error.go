package ui

import (
	"time"

	"github.com/nsf/termbox-go"
)


func ShowErrorMessage(message string) {
	termbox.Clear(termbox.ColorRed, termbox.ColorDefault)
	printMessage(0, 0, termbox.ColorWhite, termbox.ColorRed, message)
	termbox.Flush()
	time.Sleep(2 * time.Second)
}

func ShowSuccessMessage(message string) {
	termbox.Clear(termbox.ColorGreen, termbox.ColorDefault)
	printMessage(0, 0, termbox.ColorWhite, termbox.ColorGreen, message)
	termbox.Flush()
	time.Sleep(2 * time.Second)
}

func ShowTransientMessage(message string, duration time.Duration) {
	termbox.Clear(termbox.ColorYellow, termbox.ColorDefault)
	printMessage(0, 0, termbox.ColorBlack, termbox.ColorYellow, message)
	termbox.Flush()
	time.Sleep(duration)
}