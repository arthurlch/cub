package state

type Snapshot struct {
	TextBuffer [][]rune
	CurrentRow int
	CurrentCol int
}

func Capture(s *State) Snapshot {
	return Snapshot{
		TextBuffer: DeepCopyTextBuffer(s.TextBuffer),
		CurrentRow: s.CurrentRow,
		CurrentCol: s.CurrentCol,
	}
}

func (snap Snapshot) restore(s *State) {
	s.TextBuffer = DeepCopyTextBuffer(snap.TextBuffer)
	s.CurrentRow = snap.CurrentRow
	s.CurrentCol = snap.CurrentCol
	s.ClampCursor()
}

type Command interface {
	Execute(s *State)
	Undo(s *State)
	Name() string
}

type snapshotCommand struct {
	name   string
	before Snapshot
	after  Snapshot
}

func (c snapshotCommand) Name() string     { return c.name }
func (c snapshotCommand) Execute(s *State) { c.after.restore(s) }
func (c snapshotCommand) Undo(s *State)    { c.before.restore(s) }

const maxHistoryDepth = 500

type History struct {
	undo []Command
	redo []Command
}

func (s *State) Do(name string, mutate func()) {
	before := Capture(s)
	mutate()
	after := Capture(s)
	s.history.undo = append(s.history.undo, snapshotCommand{name: name, before: before, after: after})
	if len(s.history.undo) > maxHistoryDepth {
		s.history.undo = s.history.undo[len(s.history.undo)-maxHistoryDepth:]
	}
	s.history.redo = nil
}

func (s *State) Push(c Command) {
	s.history.undo = append(s.history.undo, c)
	s.history.redo = nil
}

func (s *State) Undo() bool {
	n := len(s.history.undo)
	if n == 0 {
		return false
	}
	c := s.history.undo[n-1]
	s.history.undo = s.history.undo[:n-1]
	c.Undo(s)
	s.history.redo = append(s.history.redo, c)
	return true
}

func (s *State) Redo() bool {
	n := len(s.history.redo)
	if n == 0 {
		return false
	}
	c := s.history.redo[n-1]
	s.history.redo = s.history.redo[:n-1]
	c.Execute(s)
	s.history.undo = append(s.history.undo, c)
	return true
}

func (s *State) CanUndo() bool { return len(s.history.undo) > 0 }

func (s *State) CanRedo() bool { return len(s.history.redo) > 0 }
