package manager_api

type VertexId string

type Vertex struct {
	Id      VertexId `json:"id"`
	GeoData any      `json:"geo_data"` // TODO init with real type
}

type Edge struct {
	From VertexId `json:"from"`
	To   VertexId `json:"to"`
}
