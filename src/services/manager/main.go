package main

import (
	"fmt"
	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/badochov/distributed-shortest-path/src/services/manager/executor"
	"github.com/badochov/distributed-shortest-path/src/services/manager/server"
	"log"
	"math/bits"
	"os"
	"strconv"
)

func validateNumRegions(numRegions int) error {
	if numRegions <= 0 {
		return fmt.Errorf("number of regions must be positive")
	}
	if bits.OnesCount(uint(numRegions)) != 1 {
		return fmt.Errorf("number of regions must be a power of two")
	}
	return nil
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
	if err := validateNumRegions(numRegions); err != nil {
		log.Fatal("Error validating number of regions,", err)
	}
	executorDeps := executor.Deps{
		NumRegions:        numRegions,
		RegionUrlTemplate: os.Getenv("REGION_URL_TEMPLATE"),
		Db:                orm,
	}
	exctr := executor.New(executorDeps)

	serverDeps := server.Deps{
		Executor: exctr,
	}
	srv := server.New(serverDeps)

	if err := srv.Run(); err != nil {
		log.Fatal("Error while running server,", err)
	}
}
