package document_test

import (
	"testing"

	"github.com/arthurlch/cub/internal/document"
	"github.com/stretchr/testify/assert"
)

func TestExtension(t *testing.T) {
	assert.Equal(t, "go", document.Extension("main.go"))
	assert.Equal(t, "md", document.Extension("/docs/README.MD"))
	assert.Equal(t, "", document.Extension("Makefile"))
}

func TestHasImageMarkup(t *testing.T) {
	assert.True(t, document.HasImageMarkup("md"))
	assert.True(t, document.HasImageMarkup("markdown"))
	assert.True(t, document.HasImageMarkup("html"))

	assert.False(t, document.HasImageMarkup("go"))
	assert.False(t, document.HasImageMarkup("log"))
	assert.False(t, document.HasImageMarkup("sum"))
}

func TestIsPlainText(t *testing.T) {
	assert.True(t, document.IsPlainText("log"))
	assert.True(t, document.IsPlainText("sum"))
	assert.False(t, document.IsPlainText("go"))
}
