package model

type Greeting struct {
	ID     uint `gorm:"primaryKey"`
	Prefix string
}
