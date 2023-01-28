package testing

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMockConn(dbFilePath string) (*gorm.DB, error) {
	dialector := sqlite.Open(dbFilePath)

	return gorm.Open(dialector, &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
}
