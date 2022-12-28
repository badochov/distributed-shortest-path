package main

import (
	"fmt"
	"github.com/badochov/distributed-shortest-path/src/services/worker/service"
	"log"
	"net"
)

const workerServicePort int = 1337

func newListener(port int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf(":%d", port))
}

func main() {
	l, err := newListener(workerServicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serviceDeps := service.Deps{
		Listener: l,
	}
	s := service.New(serviceDeps)

	go func() {
		if err := s.Run(); err != nil {
			log.Fatalf("error running worker service, %s", err)
		}
	}()
}
