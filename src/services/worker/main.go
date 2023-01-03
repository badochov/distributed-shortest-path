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
	"os"
	"strconv"
	"sync"
)

func getPortFromEnv(envName string) int {
	portStr := os.Getenv(envName)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error parsing port from %s, %s", envName, err)
	}
	return port
}

func main() {
	log.Println(os.Environ())

	orm, err := db.ConnectToDefault()
	if err != nil {
		log.Fatal("Error opening db,", err)
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	clientset := kubernetes.NewForConfigOrDie(config)
	discovererDeps := discoverer.Deps{
		Client:        clientset,
		Namespace:     os.Getenv("NAMESPACE"),
		LabelSelector: os.Getenv("WORKER_SERVICE_LABEL_SELECTOR"),
	}
	d := discoverer.New(discovererDeps)

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
		Port:     getPortFromEnv("WORKER_SERVER_PORT"),
		Executor: exec,
	}
	sW := service_server.New(serviceDeps)

	lL, err := net.Listen("tcp", fmt.Sprintf(":%d", getPortFromEnv("LINK_SERVER_PORT")))
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

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		if err := sW.Run(); err != nil {
			log.Fatalf("error running worker service server, %s", err)
		}
		wg.Done()
	}()
	go func() {
		if err := sL.Run(); err != nil {
			log.Fatalf("error running link service server, %s", err)
		}
		wg.Done()
	}()

	wg.Wait()
}
