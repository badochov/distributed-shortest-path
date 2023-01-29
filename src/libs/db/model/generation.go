package model

type GenerationType uint8

const (
	GenerationCurrent GenerationType = iota
	GenerationNext
	GenerationActive
)

type Generation struct {
	Generation     uint16
	GenerationType `gorm:"primary_key"`
}
