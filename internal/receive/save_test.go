package receive

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/tech-dipesh/drop/internal/wire"
)

func TestValidName(t *testing.T) {
	for _, name := range []string{"photo.jpg", "notes 2.pdf", ".bashrc"} {
		if !validName(name) {
			t.Errorf("%q should be allowed", name)
		}
	}
	for _, name := range []string{"", ".", "..", "../x", "a/b", `a\b`, "/etc/passwd", "c:x"} {
		if validName(name) {
			t.Errorf("%q should be rejected", name)
		}
	}
}

func TestFreePath(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.txt", "a (1).txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if got, want := freePath(dir, "a.txt"), filepath.Join(dir, "a (2).txt"); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestSave(t *testing.T) {
	dir := t.TempDir()
	path, err := save(dir, wire.File{Name: "x.bin", Size: 5}, bytes.NewReader([]byte("hello")), func(int64) {})
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil || string(got) != "hello" {
		t.Fatalf("got %q, %v", got, err)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 1 {
		t.Fatalf("expected only the saved file, found %d entries", len(entries))
	}
}

func TestSaveCleansUpOnFailure(t *testing.T) {
	dir := t.TempDir()
	_, err := save(dir, wire.File{Name: "x.bin", Size: 10}, bytes.NewReader([]byte("abc")), func(int64) {})
	if err == nil {
		t.Fatal("expected an error for a short stream")
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Fatalf("left %d files behind", len(entries))
	}
}
