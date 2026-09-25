package filetree_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arthurlch/cub/internal/filetree"
	"github.com/stretchr/testify/assert"
)

func mk(t *testing.T, path string, dir bool) {
	t.Helper()
	var err error
	if dir {
		err = os.Mkdir(path, 0o755)
	} else {
		err = os.WriteFile(path, []byte("x"), 0o644)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewSortsDirsFirst(t *testing.T) {
	dir := t.TempDir()
	mk(t, filepath.Join(dir, "b.txt"), false)
	mk(t, filepath.Join(dir, "a.txt"), false)
	mk(t, filepath.Join(dir, "zsub"), true)

	root, err := filetree.New(dir)
	assert.NoError(t, err)

	nodes := filetree.Flatten(root)
	assert.Equal(t, 3, len(nodes))
	assert.Equal(t, "zsub", nodes[0].Name)
	assert.True(t, nodes[0].IsDir)
	assert.Equal(t, "a.txt", nodes[1].Name)
	assert.Equal(t, "b.txt", nodes[2].Name)
}

func TestToggleExpandsAndCollapses(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	mk(t, sub, true)
	mk(t, filepath.Join(sub, "inner.txt"), false)

	root, err := filetree.New(dir)
	assert.NoError(t, err)

	nodes := filetree.Flatten(root)
	assert.Equal(t, 1, len(nodes))
	assert.False(t, nodes[0].Expanded)

	nodes[0].Toggle()
	nodes = filetree.Flatten(root)
	assert.Equal(t, 2, len(nodes))
	assert.Equal(t, "inner.txt", nodes[1].Name)
	assert.Equal(t, 2, nodes[1].Depth)

	nodes[0].Toggle()
	assert.Equal(t, 1, len(filetree.Flatten(root)))
}
