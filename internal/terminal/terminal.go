package terminal

import (
	"os"
	"os/exec"
	"sync"

	"github.com/charmbracelet/x/vt"
	"github.com/creack/pty"
)

// PTY session that runs a shell and emulates a terminal.
// It works well but depending on initial user config for terminal its sometime quite
// slow to start up.

type Session struct {
	mu      sync.Mutex
	emu     *vt.Emulator
	ptmx    *os.File
	cmd     *exec.Cmd
	updates chan struct{}
	closed  bool
}

func New(cols, rows int) (*Session, error) {
	cols, rows = clampSize(cols, rows)
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}
	cmd := exec.Command(shell)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
	if err != nil {
		return nil, err
	}

	s := &Session{
		emu:     vt.NewEmulator(cols, rows),
		ptmx:    ptmx,
		cmd:     cmd,
		updates: make(chan struct{}, 1),
	}
	go s.readLoop()
	return s, nil
}

func (s *Session) readLoop() {
	buf := make([]byte, 32768)
	for {
		n, err := s.ptmx.Read(buf)
		if n > 0 {
			s.mu.Lock()
			_, _ = s.emu.Write(buf[:n])
			s.mu.Unlock()
			s.notify()
		}
		if err != nil {
			s.mu.Lock()
			s.closed = true
			s.mu.Unlock()
			s.notify()
			return
		}
	}
}

func (s *Session) Updates() <-chan struct{} { return s.updates }

func (s *Session) notify() {
	select {
	case s.updates <- struct{}{}:
	default:
	}
}

func (s *Session) Write(p []byte) {
	if len(p) == 0 || s.Closed() {
		return
	}
	_, _ = s.ptmx.Write(p)
}

func (s *Session) Render() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.emu.Render()
}

func (s *Session) Resize(cols, rows int) {
	cols, rows = clampSize(cols, rows)
	s.mu.Lock()
	s.emu.Resize(cols, rows)
	s.mu.Unlock()
	_ = pty.Setsize(s.ptmx, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
}

func (s *Session) Closed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}

func (s *Session) Close() {
	s.mu.Lock()
	s.closed = true
	s.mu.Unlock()
	if s.cmd != nil && s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
	}
	_ = s.ptmx.Close()
}

func clampSize(cols, rows int) (int, int) {
	if cols < 1 {
		cols = 1
	}
	if rows < 1 {
		rows = 1
	}
	return cols, rows
}
