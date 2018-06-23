package file

import (
	"os"
	"log"
	"strings"
	"github.com/peter-mueller/eventstore"
	"bufio"
	"errors"
	"encoding/json"
)

type (
	Store struct {
		path        string
		writer      *os.File
		Information eventstore.Information
	}

	Options struct {
		Information eventstore.Information
	}
)

var (
	ErrDataContainsNewLine    = errors.New("cannot be saved, data contains a new line rune")
	ErrNoInformation          = errors.New("no store information present in file")
	ErrNotMatchingInformation = errors.New("present store information does not match the desired one")
)

func (s *Store) Push(data []byte) error {
	str := string(data)
	if strings.ContainsRune(str, '\n') {
		return ErrDataContainsNewLine
	}
	s.writer.Write(append(data, byte('\n')))
	return nil
}

func (s *Store) StreamAll() (<-chan eventstore.Event, error) {
	c := make(chan eventstore.Event, 23)

	file, err := os.OpenFile(s.path, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return c, err
	}
	go func() {
		defer file.Close()
		defer close(c)
		scanner := bufio.NewScanner(file)

		i := s.Information.StartIndex
		scanner.Scan() // skip later used info block //TODO
		i++

		for scanner.Scan() {
			c <- eventstore.Event{
				Data:  scanner.Bytes(),
				Index: i - 1,
			}
			i++
		}
		if err := scanner.Err(); err != nil {
			log.Printf("error happened at reading eventstore after file line %s: %s", i, err)
		}
	}()

	return c, nil
}

func (s *Store) readInformation() (eventstore.Information, error) {
	file, err := os.OpenFile(s.path, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return eventstore.Information{}, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return eventstore.Information{}, err
	}
	if stat.Size() == 0 {
		return eventstore.Information{}, ErrNoInformation
	}

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return eventstore.Information{}, err
	}
	informationLine := scanner.Bytes()
	var information eventstore.Information
	err = json.Unmarshal(informationLine, &information)
	return information, err
}

func OpenFileStore(file string, options Options) (*Store, error) {
	store := &Store{
		path: file,
	}

	writer, err := os.OpenFile(store.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	store.writer = writer

	info, err := store.readInformation()
	if err == ErrNoInformation {
		data, err := json.Marshal(options.Information)
		if err != nil {
			return store, err
		}
		store.Push(data)
		return store, nil
	}
	if info != options.Information {
		return store, ErrNotMatchingInformation
	}
	return store, err
}

func (s *Store) Close() error {
	return s.writer.Close()
}
