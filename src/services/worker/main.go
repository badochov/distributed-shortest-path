package main

import (
	"fmt"
	"github.com/badochov/distributed-shortest-path/src/services/worker/link/link_server"
	"github.com/badochov/distributed-shortest-path/src/services/worker/service_server"
	"log"
	"net"
)

const workerServicePort int = 1337
const linkServicePort int = 4567

func newListener(port int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf(":%d", port))
}

func main() {
	lW, err := newListener(workerServicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serviceDeps := service_server.Deps{
		Listener: lW,
	}
	sW := service_server.New(serviceDeps)

	go func() {
		if err := sW.Run(); err != nil {
			log.Fatalf("error running worker service server, %s", err)
		}
	}()

	lL, err := newListener(linkServicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	linkDeps := link_server.Deps{
		Listener: lL,
	}
	sL := link_server.New(linkDeps)

	go func() {
		if err := sL.Run(); err != nil {
			log.Fatalf("error running link service server, %s", err)
		}
	}()
}
