package api

const AddEdgesUrl = "/add_edges"

type AddEdgesRequest struct {
	Edges []Edge `json:"edges"`
}

type AddEdgesResponse struct {
}
