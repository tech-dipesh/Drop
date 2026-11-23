package receive

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tech-dipesh/drop/internal/transfer"
	"github.com/tech-dipesh/drop/internal/wire"
)

func validName(name string) bool {
	return name != "." && filepath.IsLocal(name) && !strings.ContainsAny(name, `/\:`)
}

func save(dir string, f wire.File, src io.Reader, add func(int64)) (string, error) {
	part, err := os.CreateTemp(dir, f.Name+".*.part")
	if err != nil {
		return "", err
	}
	if err := errors.Join(transfer.Copy(part, src, f.Size, add), part.Chmod(0o644), part.Close()); err != nil {
		os.Remove(part.Name())
		return "", err
	}
	path := freePath(dir, f.Name)
	return path, os.Rename(part.Name(), path)
}

func freePath(dir, name string) string {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	path := filepath.Join(dir, name)
	for i := 1; ; i++ {
		if _, err := os.Lstat(path); err != nil {
			return path
		}
		path = filepath.Join(dir, fmt.Sprintf("%s (%d)%s", base, i, ext))
	}
}
