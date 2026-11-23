package wire

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"
)

func TestFrameRoundTrip(t *testing.T) {
	msgs := []string{"hello", "", "second frame"}
	var buf bytes.Buffer
	for _, msg := range msgs {
		if err := writeFrame(&buf, []byte(msg)); err != nil {
			t.Fatal(err)
		}
	}
	for _, want := range msgs {
		got, err := readFrame(&buf)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	}
}

func TestFrameTooBig(t *testing.T) {
	err := writeFrame(io.Discard, make([]byte, maxFrame+1))
	if !errors.Is(err, errTooBig) {
		t.Fatalf("write got %v, want errTooBig", err)
	}
	head := binary.BigEndian.AppendUint32(nil, maxFrame+1)
	_, err = readFrame(bytes.NewReader(head))
	if !errors.Is(err, errTooBig) {
		t.Fatalf("read got %v, want errTooBig", err)
	}
}

func TestFrameCutShort(t *testing.T) {
	head := binary.BigEndian.AppendUint32(nil, 10)
	_, err := readFrame(bytes.NewReader(append(head, "abc"...)))
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("got %v, want unexpected EOF", err)
	}
}
