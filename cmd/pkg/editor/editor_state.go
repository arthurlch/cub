package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
)

type EditorState struct {
	State *state.State
}

func NewEditorState(sharedState *state.State) *EditorState {
	if sharedState == nil {
		sharedState = &state.State{}
	}
	return &EditorState{State: sharedState}
}
