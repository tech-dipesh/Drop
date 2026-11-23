package transfer

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestCopyExactSize(t *testing.T) {
	size := int64(2*chunk + 123)
	first := bytes.Repeat([]byte("a"), int(size))
	next := []byte("next file")
	src := bytes.NewReader(append(first, next...))
	var dst bytes.Buffer
	var counted int64
	calls := 0
	add := func(n int64) {
		counted += n
		calls++
	}
	if err := Copy(&dst, src, size, add); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(dst.Bytes(), first) {
		t.Fatal("copied bytes don't match")
	}
	if counted != size || calls != 3 {
		t.Fatalf("counted %d in %d calls, want %d in 3", counted, calls, size)
	}
	rest, _ := io.ReadAll(src)
	if !bytes.Equal(rest, next) {
		t.Fatalf("next file got %q, want %q", rest, next)
	}
}

func TestCopyShortSource(t *testing.T) {
	var counted int64
	err := Copy(io.Discard, bytes.NewReader([]byte("abc")), 10, func(n int64) { counted += n })
	if !errors.Is(err, io.EOF) {
		t.Fatalf("got %v, want EOF", err)
	}
	if counted != 3 {
		t.Fatalf("counted %d, want 3", counted)
	}
}
