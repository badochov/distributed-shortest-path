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

	vertices, edges, err := parser.Parse(*path)
	if err != nil {
		log.Fatalf("Error parsing file, %s: %s", *path, err)
	}

	c := client.New(http.DefaultClient, *addr)
	if err := c.AddVertices(vertices); err != nil {
		log.Fatalf("Error adding vertices, %s", err)
	}
	if err := c.AddEdges(edges); err != nil {
		log.Fatalf("Error adding edges, %s", err)
	}
}
