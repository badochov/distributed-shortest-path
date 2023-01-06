package manager_api

const AddVerticesUrl = "/add_vertices"

type AddVerticesRequest struct {
	Vertices []Vertex `json:"vertices"`
}

type AddVerticesResponse struct {
}
