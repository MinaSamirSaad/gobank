package main

import (
	"log"
)

func main() {
	store, err := PostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIserver(":3000", store)
	server.Run()

}
