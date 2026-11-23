package wire

import (
	"encoding/json"
	"fmt"
	"io"
)

const Version = 1

type Message interface{ message() }

type File struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type Offer struct {
	Version int    `json:"version"`
	From    string `json:"from"`
	Files   []File `json:"files"`
}

type Answer struct {
	Accept bool `json:"accept"`
}

type Trailer struct {
	Hashes []string `json:"hashes"`
}

type Result struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

func (Offer) message()   {}
func (Answer) message()  {}
func (Trailer) message() {}
func (Result) message()  {}

func (o Offer) Total() int64 {
	var total int64
	for _, f := range o.Files {
		total += f.Size
	}
	return total
}

func Send(w io.Writer, m Message) error {
	b, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("encode %T: %w", m, err)
	}
	return writeFrame(w, b)
}

func Receive(r io.Reader, m Message) error {
	b, err := readFrame(r)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, m); err != nil {
		return fmt.Errorf("decode %T: %w", m, err)
	}
	return nil
}
