package state

import "time"

const messageTTL = 4 * time.Second

func (s *State) ShowMessage(message string) {
	s.ErrorMessage = message
	s.MessageTimestamp = time.Now()
}

func (s *State) ActiveMessage() (string, bool) {
	if s.ErrorMessage == "" {
		return "", false
	}
	if time.Since(s.MessageTimestamp) > messageTTL {
		return "", false
	}
	return s.ErrorMessage, true
}
