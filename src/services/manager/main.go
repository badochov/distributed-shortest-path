package main

import (
	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/badochov/distributed-shortest-path/src/services/manager/executor"
	"github.com/badochov/distributed-shortest-path/src/services/manager/server"
	"log"
	"os"
	"strconv"
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

	numRegionsStr := os.Getenv("NUM_REGIONS")
	numRegions, err := strconv.Atoi(numRegionsStr)
	if err != nil {
		log.Fatal("Error parsing number of regions,", err)
	}

	executorDeps := executor.Deps{
		NumRegions:        numRegions,
		RegionUrlTemplate: os.Getenv("REGION_URL_TEMPLATE"),
		Port:              getPortFromEnv("WORKER_SERVER_PORT"),
		Db:                orm,
	}
	exctr := executor.New(executorDeps)

	serverDeps := server.Deps{
		Executor: exctr,
		Port:     getPortFromEnv("PORT"),
	}
	srv := server.New(serverDeps)

	if err := srv.Run(); err != nil {
		log.Fatal("Error while running server,", err)
	}
}
