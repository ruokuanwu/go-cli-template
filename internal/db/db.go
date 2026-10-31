package db

import (
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open initializes GORM with the supplied driver and returns its cleanup function.
func Open(dialector gorm.Dialector, options ...gorm.Option) (*gorm.DB, func() error, error) {
	options = append([]gorm.Option{&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}}, options...)
	database, err := gorm.Open(dialector, options...)
	if err != nil {
		return nil, nil, oops.Wrapf(err, "open database")
	}
	sqlDB, err := database.DB()
	if err != nil {
		return nil, nil, oops.Wrapf(err, "get database connection")
	}
	return database, sqlDB.Close, nil
}
