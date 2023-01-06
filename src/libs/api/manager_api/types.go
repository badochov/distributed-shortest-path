package manager_api

type VertexId = int64

type Vertex struct {
	Id        VertexId `json:"id"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
}

type Edge struct {
	From VertexId `json:"from"`
	To   VertexId `json:"to"`
}
