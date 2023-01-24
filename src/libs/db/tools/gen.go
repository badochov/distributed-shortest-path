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

	g.ApplyBasic(model.List...)

	g.Execute()
}
