package main

import (
	"context"
	"fmt"
	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/badochov/distributed-shortest-path/src/services/worker/discoverer"
	"github.com/badochov/distributed-shortest-path/src/services/worker/link/link_server"
	"github.com/badochov/distributed-shortest-path/src/services/worker/service_server"
	"github.com/badochov/distributed-shortest-path/src/services/worker/service_server/executor"
	"github.com/badochov/distributed-shortest-path/src/services/worker/worker"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"net"
)

const workerServicePort int = 8080
const linkServicePort int = 4567

func newListener(port int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf(":%d", port))
}

func main() {
	orm, err := db.ConnectToDefault()
	if err != nil {
		log.Fatal("Error opening db,", err)
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	clientset := kubernetes.NewForConfigOrDie(config)
	discovererDeps := discoverer.Deps{Client: clientset}
	d := discoverer.New(discovererDeps) // TODO

	workerDeps := worker.Deps{
		Db:         orm,
		Discoverer: d,
	}
	wrkr := worker.New(workerDeps)

	execDeps := executor.Deps{
		Worker: wrkr,
	}
	exec := executor.New(execDeps)

	serviceDeps := service_server.Deps{
		Port:     workerServicePort,
		Executor: exec,
	}
	sW := service_server.New(serviceDeps)

	lL, err := newListener(linkServicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	linkDeps := link_server.Deps{
		Listener: lL,
		Worker:   wrkr,
	}
	sL := link_server.New(linkDeps)

	if err := wrkr.Run(context.Background()); err != nil {
		log.Fatalf("error running worker, %s", err)
	}
	go func() {
		if err := sW.Run(); err != nil {
			log.Fatalf("error running worker service server, %s", err)
		}
	}()
	go func() {
		if err := sL.Run(); err != nil {
			log.Fatalf("error running link service server, %s", err)
		}
	}()
}
