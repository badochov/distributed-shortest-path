package api

type ShortestPathRequest struct {
	From VertexId `json:"from"`
	To   VertexId `json:"to"`
}

type ShortestPathResponse struct {
	Distance int        `json:"distance"`
	Vertices []VertexId `json:"vertices"`
}
