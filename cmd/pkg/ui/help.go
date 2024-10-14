package ui

import (
	"github.com/nsf/termbox-go"
)

func ShowHelpModal() {
	termbox.Clear(termbox.ColorDefault, SoftBlack)

	asciiArt := []string{
		" ______     __  __     ______    ",
		"/\\  ___\\   /\\ \\/\\ \\   /\\  == \\   ",
		"\\ \\ \\____  \\ \\ \\_\\ \\  \\ \\  __<   ",
		" \\ \\_____\\  \\ \\_____\\  \\ \\_____\\ ",
		"  \\/_____/   \\/_____/   \\/_____/ ",
	}

	shortcuts := []string{
		"Ctrl + H  : Show help",
		"Ctrl + Q  : Quit the editor",
		"Ctrl + U  : Undo the last change",
		"Ctrl + R  : Redo the last undo",
		"Ctrl + S  : Save the file",
		"i         : Switch to insert mode",
		"s         : Start selection",
		"z         : End selection",
		"x         : Cut selection",
		"c         : Copy selection",
		"v         : Paste selection",
		"Arrow Keys: Navigation",
		"o, p, k, l: Navigation (alternative keys)",
		"",
		"Tap any key to exit",
	}

	width, height := termbox.Size()

	totalHeight := len(asciiArt) + len(shortcuts) + 2
	startY := (height / 2) - (totalHeight / 2)

	centerText := func(text string, y int, fgColor termbox.Attribute) {
		x := (width / 2) - (len(text) / 2)
		for i, ch := range text {
			if x+i < width && y < height {
				termbox.SetCell(x+i, y, ch, fgColor, SoftBlack)
			}
		}
	}

	for i, line := range asciiArt {
		centerText(line, startY+i, ModalTextColor)
	}

	for i, shortcut := range shortcuts {
		centerText(shortcut, startY+len(asciiArt)+2+i, ModalTextColor)
	}

	termbox.Flush()
	termbox.PollEvent()
	termbox.Clear(termbox.ColorDefault, SoftBlack)
	termbox.Flush()
}
