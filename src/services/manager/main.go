package main

import (
	"context"
	"fmt"
	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/badochov/distributed-shortest-path/src/services/manager/executor"
	"github.com/badochov/distributed-shortest-path/src/services/manager/server"
	"github.com/badochov/distributed-shortest-path/src/services/manager/worker/service_manager"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"strconv"
	"time"
)

func getPortFromEnv(envName string) int {
	portStr := os.Getenv(envName)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error parsing port from %s, %s", envName, err)
	}
	return port
}

func getServer(ctx context.Context) (server.Server, error) {
	orm, err := db.ConnectToDefault()
	if err != nil {
		return nil, fmt.Errorf("error opening db, %w", err)
	}

	numRegionsStr := os.Getenv("NUM_REGIONS")
	numRegions, err := strconv.Atoi(numRegionsStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing number of regions, %w", err)
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset := kubernetes.NewForConfigOrDie(config)
	workerServerManagerDeps := service_manager.Deps{
		Client:                   clientset,
		Namespace:                os.Getenv("NAMESPACE"),
		NumRegions:               numRegions,
		WorkerDeploymentTemplate: os.Getenv("REGION_DEPLOYMENT_TEMPLATE"),
	}
	workerServerManager := service_manager.New(workerServerManagerDeps)

	defaultWorkerReplicasStr := os.Getenv("DEFAULT_WORKER_REPLICAS")
	defaultWorkerReplicas, err := strconv.ParseInt(defaultWorkerReplicasStr, 10, 32)
	executorDeps := executor.Deps{
		NumRegions:            numRegions,
		RegionUrlTemplate:     os.Getenv("REGION_URL_TEMPLATE"),
		Port:                  getPortFromEnv("WORKER_SERVER_PORT"),
		Db:                    orm,
		WorkerServerManager:   workerServerManager,
		DefaultWorkerReplicas: int32(defaultWorkerReplicas),
	}
	exctr, err := executor.New(ctx, executorDeps)
	if err != nil {
		return nil, fmt.Errorf("error starting executor, %w", err)
	}

	serverDeps := server.Deps{
		Executor: exctr,
		Port:     getPortFromEnv("PORT"),
	}
	srv := server.New(serverDeps)

	return srv, nil
}

func main() {
	ctx, can := context.WithTimeout(context.Background(), 15*time.Second)
	srv, err := getServer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	can()

	if err := srv.Run(); err != nil {
		log.Fatal("Error while running server,", err)
	}
}
