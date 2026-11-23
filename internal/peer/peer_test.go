package peer

import "testing"

func TestLoadCreatesThenKeepsIdentity(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("AppData", dir)
	t.Setenv("HOME", dir)
	first, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if first.ID == "" || first.Name == "" {
		t.Fatalf("got an empty identity %+v", first)
	}
	first.Name = "test laptop"
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
