package model

type ArcFlag struct {
	ID         uint `gorm:"primary_key,auto_increment"`
	EdgeId     int64
	Edge       Edge
	Flag       uint64
	Generation uint16
}
