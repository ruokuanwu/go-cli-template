package repository

import (
	"example.com/go-cli-template/internal/model"
	"github.com/samber/oops"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Greeting struct {
	db *gorm.DB
}

func NewGreeting(db *gorm.DB) *Greeting {
	return &Greeting{db: db}
}

func (r *Greeting) Init() error {
	if err := r.db.AutoMigrate(&model.Greeting{}); err != nil {
		return oops.Wrapf(err, "migrate greetings")
	}
	// Keep the default row stable without overwriting persisted changes.
	greeting := model.Greeting{ID: 1, Prefix: "Hello"}
	return oops.Wrapf(r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&greeting).Error, "seed greeting")
}

func (r *Greeting) Get() (model.Greeting, error) {
	var greeting model.Greeting
	err := r.db.First(&greeting, 1).Error
	return greeting, oops.Wrapf(err, "get greeting")
}
