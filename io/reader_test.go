package io

import (
	"bytes"
	"fmt"
	"testing"
)

func TestEventReader_StreamAll(t *testing.T) {
	data := bytes.NewBuffer([]byte(
		"1\n23\n4\n5\n6\n"))
	expected := []string{
		"1", "23", "4", "5", "6",
	}

	c, cerr := StreamAll(data)

	var amount int
	for event := range c {
		amount++
		if string(event.Data) != expected[event.Index] {
			t.Errorf("expected event %s, got %s", expected[event.Index], string(event.Data))
		}
	}
	if amount != len(expected) {
		t.Errorf("expected amount of %d, but got %d", amount, len(expected))
	}
	if err := <-cerr; err != nil {
		t.Errorf("got error in reader: %s", err)
	}
}

func TestEventReader_StreamAllEmpty(t *testing.T) {
	data := bytes.NewBuffer([]byte{})

	c, _ := StreamAll(data)
	for range c {
		t.Error("stream all returned one event, but should not")
	}
}

type ErrReader struct {
}

func (ErrReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("default error from ErrReader")
}

func TestEventReader_StreamAllError(t *testing.T) {
	c, cerr := StreamAll(ErrReader{})

	select {
	case <-cerr:
		// ok got error
		return
	case <-c:
		t.Errorf("received unexpected event")
	}
	t.Errorf("did not exit from error")
}
