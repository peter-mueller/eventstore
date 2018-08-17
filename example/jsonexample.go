package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/peter-mueller/eventstore/filestore"
)

// HadCoffeeWith as example event
type HadCoffeeWith struct {
	Name string
	Date time.Time
}

func main() {
	store, err := filestore.OpenFileStore("eventlog")
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	date := HadCoffeeWith{Date: time.Now(), Name: "Anna"}
	data, _ := json.Marshal(date)

	// push one date
	store.Push(data)

	// log the dates
	c, _ := store.StreamAll()
	for data := range c {
		var event HadCoffeeWith
		err := json.Unmarshal(data.Data, &event)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Had a date with %s at %s on %s.\n",
			event.Name,
			event.Date.Format("15:04"),
			event.Date.Weekday(),
		)
	}
	if store.LastErr != nil {
		log.Fatal(store.LastErr)
	}
}
