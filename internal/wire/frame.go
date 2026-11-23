package wire

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const maxFrame = 1 << 20

var errTooBig = errors.New("frame too big")

func writeFrame(w io.Writer, b []byte) error {
	if len(b) > maxFrame {
		return fmt.Errorf("%w: %d bytes", errTooBig, len(b))
	}
	head := binary.BigEndian.AppendUint32(nil, uint32(len(b)))
	_, err := w.Write(append(head, b...))
	return err
}

func readFrame(r io.Reader) ([]byte, error) {
	var head [4]byte
	if _, err := io.ReadFull(r, head[:]); err != nil {
		return nil, err
	}
	n := binary.BigEndian.Uint32(head[:])
	if n > maxFrame {
		return nil, fmt.Errorf("%w: %d bytes", errTooBig, n)
	}
	b := make([]byte, n)
	_, err := io.ReadFull(r, b)
	return b, err
}
