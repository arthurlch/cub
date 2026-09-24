package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRoundtrip(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	want := Config{
		Theme:        "gruvbox",
		SidebarRight: true,
		Keys:         map[string]string{"quit": "ctrl+x"},
	}
	if err := Save(want); err != nil {
		t.Fatalf("save: %v", err)
	}
	if Path() != filepath.Join(dir, "cub", "config.json") {
		t.Fatalf("unexpected path %q", Path())
	}

	got := Load()
	if got.Theme != want.Theme || got.SidebarRight != want.SidebarRight {
		t.Errorf("got %+v, want %+v", got, want)
	}
	if got.Keys["quit"] != "ctrl+x" {
		t.Errorf("keys not persisted: %+v", got.Keys)
	}
}

func TestLoadFirstRunWritesTemplate(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfg := Load()
	if len(cfg.Keys) == 0 {
		t.Fatalf("first-run config should carry default chords")
	}
	if cfg.Keys["quit"] != "ctrl+q" {
		t.Errorf("default quit chord missing: %+v", cfg.Keys)
	}
	reload := Load()
	if reload.Keys["quit"] != "ctrl+q" {
		t.Errorf("template not persisted on first run")
	}
}

func TestOldConfigMigratesKeysButKeepsTheme(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	old := `{"theme":"gruvbox","sidebar_right":true,"keys":{"undo":"ctrl+u","select_start":"s"}}`
	path := Path()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := Load()
	if cfg.Theme != "gruvbox" || !cfg.SidebarRight {
		t.Errorf("migration must preserve theme/sidebar: %+v", cfg)
	}
	if cfg.Keys["undo"] != "u" {
		t.Errorf("stale keymap should refresh to new defaults, got undo=%q", cfg.Keys["undo"])
	}
	if _, ok := cfg.Keys["select_start"]; ok {
		t.Error("removed actions should not survive migration")
	}
}

func TestChordsMergeOverDefaults(t *testing.T) {
	c := Config{Keys: map[string]string{"quit": "ctrl+x"}}
	chords := c.Chords()
	if chords["quit"] != "ctrl+x" {
		t.Errorf("override lost: %q", chords["quit"])
	}
	if chords["save"] != "ctrl+s" {
		t.Errorf("default save chord dropped after merge: %q", chords["save"])
	}
}
