package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/badochov/distributed-shortest-path/src/cmd/osm_map_importer/client"
	"github.com/badochov/distributed-shortest-path/src/cmd/osm_map_importer/parser"
)

func main() {
	addr := flag.String("address", "localhost:8080", "Address of manager server.")
	path := flag.String("file", "local/data/monaco.osm", "Path to .osm file containing graph data.}")

	flag.Parse()
	log.Println("Starting")
	log.Println(*addr, *path)

	vertices, edges, err := parser.Parse(*path)
	if err != nil {
		log.Fatalf("Error parsing file, %s: %s", *path, err)
	}

	c := client.New(http.DefaultClient, *addr)
	if err := doInBatches(c.AddVertices, vertices, 1024); err != nil {
		log.Fatalf("Error adding vertices, %s", err)
	}
	if err := doInBatches(c.AddEdges, edges, 1024); err != nil {
		log.Fatalf("Error adding edges, %s", err)
	}
}

func doInBatches[T any](fn func([]T) error, data []T, batchSize int) error {
	for len(data) > 0 {
		if batchSize > len(data) {
			batchSize = len(data)
		}
		batch := data[:batchSize]
		if err := fn(batch); err != nil {
			return err
		}
		data = data[batchSize:]
	}
	return nil
}
