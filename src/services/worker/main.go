package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/badochov/distributed-shortest-path/src/services/worker/discoverer"
	"github.com/badochov/distributed-shortest-path/src/services/worker/link/link_server"
	"github.com/badochov/distributed-shortest-path/src/services/worker/service_server"
	"github.com/badochov/distributed-shortest-path/src/services/worker/service_server/executor"
	"github.com/badochov/distributed-shortest-path/src/services/worker/worker"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

	regionIdStr := os.Getenv("REGION")
	regionId, err := strconv.ParseUint(regionIdStr, 10, 16)
	if err != nil {
		log.Fatal("can't parser REGION", err)
	}
	linkPort := getPortFromEnv("LINK_SERVER_PORT")
	workerDeps := worker.Deps{
		Db:         orm,
		Discoverer: d,
		RegionID:   uint16(regionId),
		Context:    context.Background(),
		LinkPort:   linkPort,
	}
	wrkr, err := worker.New(workerDeps)
	if err != nil {
		log.Fatalf("error creating worker, %s", err)
	}

	ctx, can := context.WithTimeout(context.Background(), 3*time.Minute)
	if err := wrkr.LoadRegionData(ctx); err != nil {
		log.Fatalf("error loading region data, %s", err)
	}
	can()

	execDeps := executor.Deps{
		Worker: wrkr,
	}
	exec := executor.New(execDeps)

	serviceDeps := service_server.Deps{
		Port:     getPortFromEnv("WORKER_SERVER_PORT"),
		Executor: exec,
	}
	sW := service_server.New(serviceDeps)

	lL, err := net.Listen("tcp", fmt.Sprintf(":%d", linkPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	linkDeps := link_server.Deps{
		Listener: lL,
		Worker:   wrkr,
	}
	sL := link_server.New(linkDeps)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		if err := sW.Run(); err != nil {
			log.Fatalf("error running worker service server, %s", err)
		}
	}()
	go func() {
		defer wg.Done()

		if err := sL.Run(); err != nil {
			log.Fatalf("error running link service server, %s", err)
		}
	}()

	wg.Wait()
}
