package filestore

import (
	"fmt"
	"io"
	"os"

	"github.com/peter-mueller/eventstore"
	eventio "github.com/peter-mueller/eventstore/io"
)

type (
	// File is an event store in file format
	File struct {
		Path    string
		LastErr error

		writer io.WriteCloser
	}
)

var (
	//ErrNoWriter if the store was not initialised with a file to push to
	ErrNoWriter = fmt.Errorf("no writer is present in the File event store")
)

// OpenFileStore opens a file that can also be pushed to.
//
// The file will be created if it does not already exist.
func OpenFileStore(path string) (*File, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &File{
		Path:   path,
		writer: f,
	}, nil
}

// Close the file still open for pushing
func (f *File) Close() {
	if f.writer == nil {
		panic(ErrNoWriter)
	}
	f.writer.Close()
}

// Push data to the file.
func (f *File) Push(data []byte) error {
	if f.writer == nil {
		return ErrNoWriter
	}
	return eventio.Push(f.writer, data)
}

// StreamAll events from the file.
func (f *File) StreamAll() (<-chan eventstore.Event, error) {
	reader, err := os.OpenFile(f.Path, os.O_SYNC|os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	c, cerr := eventio.StreamAll(reader)
	go func() {
		for err := range cerr {
			f.LastErr = err
		}
	}()
	return c, nil
}
