package main

import (
	"context"
	"log"
	"time"

	"github.com/badochov/distributed-shortest-path/src/libs/db/cleaner"
)

func main() {
	clnr, err := cleaner.ConnectToDefault()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := clnr.Clean(ctx); err != nil {
		log.Fatal(err)
	}
}
