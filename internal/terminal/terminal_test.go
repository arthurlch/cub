package terminal

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
)

func TestSessionRunsShellCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("skips spawning a shell in short mode")
	}
	s, err := New(80, 24)
	if err != nil {
		t.Skipf("no pty available: %v", err)
	}
	defer s.Close()

	s.Write([]byte("printf CUBMARK\r"))

	deadline := time.After(5 * time.Second)
	for {
		select {
		case <-s.Updates():
			if strings.Contains(ansi.Strip(s.Render()), "CUBMARK") {
				return
			}
		case <-deadline:
			t.Fatalf("shell output never showed the marker; got:\n%s", ansi.Strip(s.Render()))
		}
	}
}

func TestSessionResizeDoesNotPanic(t *testing.T) {
	if testing.Short() {
		t.Skip("skips spawning a shell in short mode")
	}
	s, err := New(40, 10)
	if err != nil {
		t.Skipf("no pty available: %v", err)
	}
	defer s.Close()
	s.Resize(100, 30)
	s.Resize(0, 0)
	_ = s.Render()
}
