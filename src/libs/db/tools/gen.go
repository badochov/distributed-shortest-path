package main

import (
	"github.com/badochov/distributed-shortest-path/src/libs/db/model"
	"gorm.io/gen"
)

//go:generate go run gen.go

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../query",
	})

	g.ApplyBasic(model.Vertex{}, model.Edge{}, model.ArcFlag{}, model.RegionBinding{}, model.Generation{})

	g.Execute()
}
