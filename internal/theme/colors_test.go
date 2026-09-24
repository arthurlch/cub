package theme_test

import (
	"testing"

	"github.com/alecthomas/chroma/v2"
	"github.com/arthurlch/cub/internal/theme"
	"github.com/stretchr/testify/assert"
)

func TestThemesDeriveValidColors(t *testing.T) {
	assert.NotEmpty(t, theme.Themes)
	for _, th := range theme.Themes {
		assert.NotEmptyf(t, string(th.Bg), "%s bg", th.Name)
		assert.NotEmptyf(t, string(th.Fg), "%s fg", th.Name)
		assert.NotEmptyf(t, string(th.Accent), "%s accent", th.Name)
		assert.NotEmptyf(t, string(th.Gutter), "%s gutter", th.Name)
		assert.NotEmptyf(t, string(th.Token(chroma.Keyword)), "%s keyword token", th.Name)
	}
}

func TestTokenFallsBackToForeground(t *testing.T) {
	th := theme.Default()
	assert.Equal(t, th.Fg, th.Token(chroma.Text), "unset token uses the foreground")
}
