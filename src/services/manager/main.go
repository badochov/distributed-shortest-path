package main

import (
	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/badochov/distributed-shortest-path/src/services/manager/discoverer"
	"github.com/badochov/distributed-shortest-path/src/services/manager/executor"
	"github.com/badochov/distributed-shortest-path/src/services/manager/server"
	"log"
)

func main() {
	orm, err := db.ConnectToDefault()
	if err != nil {
		log.Fatal("Error opening db,", err)
	}

	discovererDeps := discoverer.Deps{}
	d := discoverer.New(discovererDeps)

	executorDeps := executor.Deps{
		Discoverer: d,
		Db:         orm,
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
