package file

import (
	"testing"
	"github.com/peter-mueller/eventstore"
	"os"
	"encoding/json"
	"time"
)

func TestOpenFileStore(t *testing.T) {
	store, err :=  OpenFileStore("test.db", Options{
		Information: eventstore.Information{
			Name: "test",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	info,err := store.readInformation()
	if err != nil {
		t.Fatal(err)
	}
	if info.StartIndex != 0 {
		t.Error("start index in new file store must be 0")
	}
	if info.Name != "test" {
		t.Errorf("store name should be test, but was %s", info.Name)
	}
}



func TestStore_Push(t *testing.T) {
	store, err :=  OpenFileStore("test.db", Options{})
	if err != nil {
		t.Fatal(err)
	}
	store.Push([]byte("hello"))
	store.Push([]byte("franz"))

	c, err := store.StreamAll()
	hello := <-c
	if hello.Index != 0 {
		t.Errorf("index for hello should be 0, but was %d", hello.Index)
	}
	if string(hello.Data ) != "hello" {
		t.Errorf("data should be hello , but was %s", string(hello.Data))
	}

	franz := <-c
	if franz.Index != 1 {
		t.Errorf("index for franz should be 1, but was %d", hello.Index)
	}
	if string(franz.Data ) != "franz" {
		t.Errorf("data should be franz , but was %s", string(hello.Data))
	}
}

func TestStore_PushFailsWithNewline(t *testing.T) {
	err := os.Remove("test.db")
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	store, err :=  OpenFileStore("test.db", Options{})
	if err != nil {
		t.Fatal(err)
	}

	err = store.Push([]byte("hello\n"))
	if err != ErrDataContainsNewLine {
		t.Fatal(err)
	}
}

func BenchmarkStore_Push(b *testing.B) {
	store, err :=  OpenFileStore("test.db", Options{})
	if err != nil {
		b.Fatal(err)
	}
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		store.Push([]byte("hello"))
	}
}

func BenchmarkStore_PushJSON(b *testing.B) {
	store, err :=  OpenFileStore("test.db", Options{})
	if err != nil {
		b.Fatal(err)
	}
	// run the Fib function b.N times
	now := time.Now()
	for n := 0; n < b.N; n++ {
		data, err := json.Marshal(struct {
			Time time.Time
			Name string
		}{
			Time: now,
			Name: "hello",
		} )
		if err != nil {
			b.Fatal(err)
		}
		store.Push(data)
	}
}