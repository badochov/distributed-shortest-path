package api

const ShortestPathUrl = "/shortest_path"

type ShortestPathRequest struct {
	From VertexId `json:"from"`
	To   VertexId `json:"to"`
}

type ShortestPathResponse struct {
	NoPath   bool       `json:"no_path,omitempty"`
	Distance float64    `json:"distance"`
	Vertices []VertexId `json:"vertices"`
}
