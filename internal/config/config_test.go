package config

import (
	"path/filepath"
	"testing"
)

func TestLoadCreatesThenKeepsSettings(t *testing.T) {
	dir := t.TempDir()
	for _, key := range []string{"XDG_CONFIG_HOME", "AppData", "HOME", "USERPROFILE"} {
		t.Setenv(key, dir)
	}
	first, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if first.ID == "" || first.Name == "" || first.Dir == "" {
		t.Fatalf("got incomplete settings %+v", first)
	}
	first.Name = "test laptop"
	first.Dir = filepath.Join(dir, "inbox")
	if err := Save(first); err != nil {
		t.Fatal(err)
	}
	again, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if again != first {
		t.Fatalf("got %+v after reload, want %+v", again, first)
	}
}
