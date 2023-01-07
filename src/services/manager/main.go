package main

import (
	"context"
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

func main() {
	orm, err := db.ConnectToDefault()
	if err != nil {
		log.Fatal("Error opening db,", err)
	}

	numRegionsStr := os.Getenv("NUM_REGIONS")
	numRegions, err := strconv.Atoi(numRegionsStr)
	if err != nil {
		log.Fatal("Error parsing number of regions,", err)
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	clientset := kubernetes.NewForConfigOrDie(config)
	workerServerManagerDeps := service_manager.Deps{
		Client:                   clientset,
		Namespace:                os.Getenv("NAMESPACE"),
		NumRegions:               numRegions,
		WorkerDeploymentTemplate: os.Getenv("REGION_DEPLOYMENT_TEMPLATE"),
	}
	workerServerManager := service_manager.New(workerServerManagerDeps)

	executorDeps := executor.Deps{
		NumRegions:          numRegions,
		RegionUrlTemplate:   os.Getenv("REGION_URL_TEMPLATE"),
		Port:                getPortFromEnv("WORKER_SERVER_PORT"),
		Db:                  orm,
		WorkerServerManager: workerServerManager,
	}
	exctr := executor.New(executorDeps)

	serverDeps := server.Deps{
		Executor: exctr,
		Port:     getPortFromEnv("PORT"),
	}
	srv := server.New(serverDeps)

	ctx, can := context.WithTimeout(context.Background(), 15*time.Second)
	defer can()
	if err := srv.Run(ctx); err != nil {
		log.Fatal("Error while running server,", err)
	}
}
