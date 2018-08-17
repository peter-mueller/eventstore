package io

import (
	"bytes"
	"fmt"
	fio "io"
)

var (
	// ErrContainsNewLine is thrown if the data contains a new line.
	ErrContainsNewLine = fmt.Errorf("data contains a new line rune, cannot be saved")
)

// Push some data to the writer.
// In data no new lines are allowed.
func Push(w fio.Writer, data []byte) error {
	if bytes.ContainsRune(data, '\n') {
		return ErrContainsNewLine
	}

	_, err := w.Write(append(data, '\n'))
	return err
}
