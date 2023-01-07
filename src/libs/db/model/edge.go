package model

type Edge struct {
	ID         int64 `gorm:"primary_key"`
	FromId     int64 `gorm:"unique_index:from_to_vertices"`
	From       Vertex
	ToId       int64 `gorm:"unique_index:from_to_vertices"`
	To         Vertex
	Generation uint16
}
