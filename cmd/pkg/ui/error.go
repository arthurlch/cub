package ui

import (
	"time"

	"github.com/arthurlch/cub/cmd/pkg/state"
)

func ShowErrorMessage(st *state.State, message string) {
	st.ErrorMessage = message
	st.MessageTimestamp = time.Now()
}

func ShowSuccessMessage(st *state.State, message string) {
	st.ErrorMessage = message
	st.MessageTimestamp = time.Now()
}

func ShowTransientMessage(st *state.State, message string, duration time.Duration) {
	st.ErrorMessage = message
	st.MessageTimestamp = time.Now().Add(duration)
}

