package transfer

import (
	"fmt"
	"io"
)

const chunk = 4 << 20

func Copy(dst io.Writer, src io.Reader, size int64, add func(int64)) error {
	for left := size; left > 0; {
		n, err := io.CopyN(dst, src, min(chunk, left))
		add(n)
		left -= n
		if err != nil {
			return fmt.Errorf("copy stopped at %d of %d bytes: %w", size-left, size, err)
		}
	}
	return nil
}
