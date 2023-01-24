package api

const ShortestPathUrl = "/shortest_path"

type RequestId uint64

type ShortestPathRequest struct {
	RequestId RequestId `json:"request_id"`
	From      VertexId  `json:"from"`
	To        VertexId  `json:"to"`
}

type ShortestPathResponse struct {
	Distance int        `json:"distance"`
	Vertices []VertexId `json:"vertices"`
}
