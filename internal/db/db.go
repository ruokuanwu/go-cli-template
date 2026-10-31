package db

import (
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open initializes GORM with a caller-provided database driver.
func Open(dialector gorm.Dialector, options ...gorm.Option) (*gorm.DB, error) {
	options = append([]gorm.Option{&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}}, options...)
	database, err := gorm.Open(dialector, options...)
	if err != nil {
		return nil, oops.Wrapf(err, "open database")
	}
	return database, nil
}
