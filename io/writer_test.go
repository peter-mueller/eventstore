package io

import (
	"bytes"
	"testing"
)

func TestPush(t *testing.T) {
	buffer := &bytes.Buffer{}

	expected := []string{
		"1", "23", "4", "5", "6",
	}
	for _, item := range expected {
		err := Push(buffer, []byte(item))
		if err != nil {
			t.Errorf("got error pushing: %s", err)
		}
	}

	c, _ := StreamAll(buffer)
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
}
