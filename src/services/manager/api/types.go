package api

type VertexId = int64
type EdgeId = int64
type EdgeLength = float64

type Vertex struct {
	Id        VertexId `json:"id" gorm:"primarykey"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
}

type Edge struct {
	From   VertexId   `json:"from"`
	To     VertexId   `json:"to"`
	Length EdgeLength `json:"length"`
	Id     EdgeId     `json:"flag_id" gorm:"primarykey"`
}
