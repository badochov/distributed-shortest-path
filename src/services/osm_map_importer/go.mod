module github.com/badochov/distributed-shortest-path/src/services/osm_map_importer

go 1.19

require (
	github.com/badochov/distributed-shortest-path/src/libs/api v0.0.0
	github.com/glaslos/go-osm v0.0.0-20170316165313-16aac6148584
)

replace github.com/badochov/distributed-shortest-path/src/libs/api v0.0.0 => ../../libs/api
