package model

type Vertex struct {
	ID         int64 `gorm:"primary_key"`
	Latitude   float64
	Longitude  float64
	Generation uint16
}
