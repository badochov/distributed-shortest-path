package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDefault() (*gorm.DB, error) {
	dsn := ""
	dialector := postgres.Open(dsn)
	return gorm.Open(dialector)
}
