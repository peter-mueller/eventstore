package io

import (
	"bufio"
	fio "io"

	"github.com/peter-mueller/eventstore"
)

// StreamAll events from the reader until it is consumed.
func StreamAll(r fio.Reader) (<-chan eventstore.Event, <-chan error) {
	c := make(chan eventstore.Event, 23)
	errc := make(chan error)
	var i uint

	scanner := bufio.NewScanner(r)
	go func() {
		for scanner.Scan() {
			c <- eventstore.Event{
				Data:  scanner.Bytes(),
				Index: i,
			}
			i++
		}
		err := scanner.Err()
		if err != nil {
			errc <- err
		}
		close(c)
		close(errc)
	}()
	return c, errc
}
